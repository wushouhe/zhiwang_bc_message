package main

import (
	"zhiwang_bc_message/geth/lostblock"
	"zhiwang_bc_message/geth/blockdb"
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/subscribe"
	"github.com/golang/glog"
	"flag"
)

func main() {
	//complete.TestBinary()
	//complete.TestArrayLoop()


	flag.Parse()
	client, _ := rpc.Dial("http://172.16.10.163:8545")
	blockChan := make(chan *json.JsonHeader, 1000)
	db := blockdb.NewDB()
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
		return
	}
	subscribe.FillBlockRangeComplete(client, blockChan, lostList)

	i := 0
	for block := range blockChan {
		i++
		if i > lostBlockLen {
			close(blockChan)
		}
		glog.Infof("%s \n", block.Number.ToInt())
		blockdb.InesrtBlockChan(db, block)
	}
}
