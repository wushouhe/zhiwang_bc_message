package main

/**
--config D:/mygo/src/zhiwang_bc_message/cfg/config.yaml  --alsologtostderr=true
*/
import (
	"fmt"
	cmdutil "github.com/17golang/golang/cmd/utils"
	"github.com/17golang/golang/goutils/config"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang/glog"
	"gopkg.in/urfave/cli.v1"
	"os"
	"sync"
	"time"
	"zhiwang_bc_message/geth/blockdb"
	. "zhiwang_bc_message/geth/config"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/lostblock"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/utils"
)

var (
	app = cmdutil.NewApp("指旺区块同步接口")
)

func init() {
	app.Action = syncBlocks
	app.Copyright = "Copyright 2017 "
}
func main() {
	app.AddFlag(utils.ConfigFileFlag)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func syncBlocks(ctx *cli.Context) error {
	if err := config.ReadConfig(ctx, utils.ConfigFileFlag.Name, &Cfg); err != nil {
		return err
	}
	cmdutil.GlogShim(ctx)
	fmt.Printf("%v \n", Cfg)
	client, _ := rpc.Dial(fmt.Sprintf("%s://%s:%s", Cfg.RPC.Protocol, Cfg.RPC.Ip, Cfg.RPC.Port))
	blockChan := make(chan *json.JsonHeader, Cfg.BlocChanSize)

	db := blockdb.NewDB()
	db.SetMaxOpenConns(Cfg.ThreadSize)
	db.SetMaxIdleConns(Cfg.ThreadSize)

	lostblock.SyncLostBlock(client, db, blockChan)
	glog.Info("开始同步区块...")
	go func() {
		for {
			subscribe.SyncBlocks(client, db, blockChan)
			time.Sleep(2 * time.Minute)
		}
	}()

	var w sync.WaitGroup
	w.Add(Cfg.ThreadSize)
	for i := 0; i < Cfg.ThreadSize; i++ {
		go func() {
			for {
				block := <-blockChan
				glog.Infof("importing block %s \n", block.Number.ToInt())
				blockdb.InesrtBlockChan(db, block)
			}
			w.Done()
		}()
	}
	w.Wait()

	return nil
}
