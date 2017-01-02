package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/json"
	//"zhiwang_bc_message/geth/utils"
	"zhiwang_bc_message/geth/blockdb"
	"fmt"
)

func main() {
	client, _ := rpc.Dial("http://139.196.178.168:8545")
	blockChan := make(chan *json.JsonHeader, 100)
	db := blockdb.NewDB()
	subscribe.SyncBlocks(client, db, blockChan)
	for block := range blockChan {
		fmt.Printf("%v \n",block)
		blockdb.InsertBlock(db, block)
	}
}
