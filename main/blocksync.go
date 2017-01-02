package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/json"
	//"zhiwang_bc_message/geth/utils"
	"zhiwang_bc_message/geth/blockdb"
	"fmt"
	"time"
)

func main() {
	client, _ := rpc.Dial("http://139.196.178.168:8545")
	blockChan := make(chan *json.JsonHeader, 100)
	db := blockdb.NewDB()
	fmt.Println("开始同步区块...")
	go subscribe.SyncBlocks(client, db, blockChan)

	timer := time.NewTimer(time.Second * 30)
	for {
		select {
		case <-timer.C:
			subscribe.SyncBlocks(client, db, blockChan)
			timer.Reset(time.Second * 30)
		case block:=<-blockChan:
			blockdb.InsertBlock(db, block)
			timer.Reset(time.Second * 30)
		}
	}
	/*for block := range blockChan {
		//fmt.Printf("%v \n", block)
		blockdb.InsertBlock(db, block)
	}*/

}
