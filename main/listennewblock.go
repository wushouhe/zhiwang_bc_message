package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/blockdb"
	"zhiwang_bc_message/geth/utils"
	"github.com/golang/glog"
	"sync"
	"zhiwang_bc_message/geth/lostblock"
	. "zhiwang_bc_message/geth/config"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"
)

var once sync.Once

var (
	app = utils.NewApp("zhiwang message midware command line interface")
)

func init() {
	app.Action = listenNewBlock
	app.Copyright = "Copyright 2017 "
	/*app.Commands = []cli.Command{
		syncBlockCommand,
		lostBlockCommand,
	}*/

	app.Flags = []cli.Flag{
		utils.ConfigFileFlag,
	}
	app.Flags = append(app.Flags, utils.GlogGangstaFlags...)

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		glog.Flush()
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func listenNewBlock(ctx *cli.Context) error {
	//解析配置文件
	err := utils.ReadConfig(ctx)
	if err != nil {
		return err
	}
	//glog
	utils.GlogGangstaShim(ctx)

	client, _ := rpc.Dial(fmt.Sprintf("%s://%s:%s", Cfg.RPC.Protocol, Cfg.RPC.Ip, Cfg.RPC.Port))

	blockChan := make(chan *json.JsonHeader, Cfg.BlocChanSize)

	glog.Infof("正在监听new Block.........")
	subscribe.ListenNewBlock(client, blockChan)

	db := blockdb.NewDB()

	lostBlockFunc := func() {
		lostblock.SyncLostBlock(client, db, blockChan)
	}

	var w sync.WaitGroup
	w.Add(Cfg.ThreadSize)
	for i := 0; i < Cfg.ThreadSize; i++ {
		go func() {
			for {
				block := <-blockChan
				glog.Infof("importing block %s \n", block.Number.ToInt())
				blockdb.InesrtBlockChan(db, block)
				once.Do(lostBlockFunc)
			}
			w.Done()
		}()
	}
	w.Wait()

	return nil
}




