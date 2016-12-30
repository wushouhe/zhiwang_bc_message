package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/subscribe"
	//"zhiwang_bc_message/geth/utils"
	"fmt"
)

func main() {
	client, _ := rpc.Dial("http://172.16.10.163:8545")
	blockChan := make(chan *json.JsonHeader,100)
	subscribe.SyncAndSubscribBlock(client, blockChan)
	for {
		select {
		case block := <-blockChan:
		//utils.PrintBlock(block)
			if block != nil {
				fmt.Printf("%v \n", block.Number)

			}else{
				fmt.Println("nil block")
			}
		}
	}

}
