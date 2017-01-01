package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/subscribe"
	//"zhiwang_bc_message/geth/utils"
	"fmt"
)

func main() {
	client, _ := rpc.Dial("http://139.196.178.168:8545")
	blockChan := make(chan *json.JsonHeader,100)
	subscribe.SyncAndSubscribBlock(client, blockChan)
	for {
		select {
		case block := <-blockChan:
			//utils.PrintBlock(block)
			fmt.Printf("%s \n",block.Number.ToInt())
		}
	}

}
