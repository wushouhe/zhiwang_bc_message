package main

import (
	cmdutil "github.com/17golang/golang/cmd/utils"
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
)

var once sync.Once

var (
	app = cmdutil.NewApp("指旺区块监听接口")
)

func init() {
	app.Action = listenNewBlock
	app.Copyright = "Copyright 2017 "
	/*app.Commands = []cli.Command{
		syncBlockCommand,
		lostBlockCommand,
	}*/

	app.AddFlag(utils.ConfigFileFlag)
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
	db.SetMaxOpenConns(Cfg.ThreadSize)
	db.SetMaxIdleConns(Cfg.ThreadSize)

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
