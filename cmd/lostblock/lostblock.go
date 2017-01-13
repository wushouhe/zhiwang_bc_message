package main

import (
	cmdutil "github.com/17golang/golang/cmd/utils"
	"zhiwang_bc_message/geth/lostblock"
	"zhiwang_bc_message/geth/blockdb"
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/subscribe"
	"github.com/golang/glog"
	"os"
	"fmt"
	. "zhiwang_bc_message/geth/config"
	"zhiwang_bc_message/geth/utils"
	"gopkg.in/urfave/cli.v1"
)

var (
	app = cmdutil.NewApp("检查丢失区块接口")
)

func init() {
	app.Action = lostBlock
	app.Copyright = "Copyright 2017 "

	app.AddFlag(utils.ConfigFileFlag)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func lostBlock(ctx *cli.Context) error {
	//解析配置文件
	err := utils.ReadConfig(ctx)
	if err != nil {
		return err
	}
	//glog
	utils.GlogGangstaShim(ctx)

	client, _ := rpc.Dial(fmt.Sprintf("%s://%s:%s", Cfg.RPC.Protocol, Cfg.RPC.Ip, Cfg.RPC.Port))

	blockChan := make(chan *json.JsonHeader, Cfg.BlocChanSize)

	db := blockdb.NewDB()
	db.SetMaxOpenConns(Cfg.ThreadSize)
	db.SetMaxIdleConns(Cfg.ThreadSize)

	//删除重复数据
	glog.Infof("删除重复区块")
	blockdb.RemoveRepeatRows(db)
	glog.Infof("检查缺失区块")
	lostList := lostblock.MysqlLoop(db)
	for k, v := range lostList {
		glog.Infof("k %d v %d ", k, v)
	}
	lostBlockLen := len(lostList)
	if lostBlockLen == 0 {
		return nil
	}
	go subscribe.FillBlockRangeComplete(client, blockChan, lostList)

	i := 0
	for block := range blockChan {

		if i >= lostBlockLen {
			close(blockChan)
			break
		}
		glog.Infof("%s \n", block.Number.ToInt())
		blockdb.InesrtBlockChan(db, block)
		i++
	}
	return nil
}
