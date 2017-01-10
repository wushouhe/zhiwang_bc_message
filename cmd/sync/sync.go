package main
/**
--config D:/mygo/src/zhiwang_bc_message/cfg/config.yaml  --alsologtostderr=true
 */
import (
	cmdutil "github.com/17golang/golang/cmd/utils"
	"gopkg.in/urfave/cli.v1"
	"github.com/17golang/golang/goutils/config"
	. "zhiwang_bc_message/geth/config"
	"zhiwang_bc_message/geth/utils"
	"os"
	"fmt"
	"sync"
	"zhiwang_bc_message/geth/blockdb"
	"zhiwang_bc_message/geth/lostblock"
	"zhiwang_bc_message/geth/subscribe"
	"time"
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"github.com/golang/glog"
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
	fmt.Printf("%v \n",Cfg)
	client, _ := rpc.Dial(fmt.Sprintf("%s://%s:%s", Cfg.RPC.Protocol, Cfg.RPC.Ip, Cfg.RPC.Port))
	blockChan := make(chan *json.JsonHeader, Cfg.BlocChanSize)
	db := blockdb.NewDB()
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