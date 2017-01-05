package main

import (
	"zhiwang_bc_message/geth/complete"
	"zhiwang_bc_message/geth/blockdb"
	//"github.com/ethereum/go-ethereum/rpc"
	//"zhiwang_bc_message/geth/json"
	//"zhiwang_bc_message/geth/subscribe"
	//"github.com/golang/glog"
	"flag"
	"fmt"
)

func main() {
	//complete.TestBinary()
	//complete.TestArrayLoop()


	flag.Parse()
	//client, _ := rpc.Dial("http://172.16.10.163:8545")
	//blockChan := make(chan *json.JsonHeader, 1000)
	db := blockdb.NewDB()
	//删除重复数据
	//blockdb.RemoveRepeatRows(db)

	list := complete.MysqlLoop(db)
	for k, v := range list {
		fmt.Printf("k %d v %d \n ", k, v)
	}

	/*subscribe.FillBlockRangeComplete(client, blockChan, list)

	for block := range blockChan {
		glog.Infof("%s \n", block.Number.ToInt())
		blockdb.InesrtBlockChan(db, block)
	}*/
}
