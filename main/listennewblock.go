package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/blockdb"
	"github.com/golang/glog"
)

func main() {
	//172.16.10.163
	//139.196.178.168
	client, _ := rpc.Dial("http://139.196.178.168:8545")
	blockChan := make(chan *json.JsonHeader, 1000)
	glog.Infof("正在监听new Block.........")
	subscribe.ListenNewBlock(client, blockChan)
	db := blockdb.NewDB()
	for {
		block := <-blockChan
		blockdb.InesrtBlockChan(db, block)
		glog.Infof("block %v ", block)
	}

}





