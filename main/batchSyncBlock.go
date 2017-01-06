package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/json"
	"fmt"
)

func main() {
	client, _ := rpc.Dial("http://172.16.10.163:8545")
	blockChan := make(chan *json.JsonHeader, 100)
	//subscribe.BatchRequest(client,blockChan,0,99)
	go subscribe.FillBlockRange(client, blockChan, 0, 706)

	for block := range blockChan {
		//utils.PrintBlock(block)
		fmt.Printf("%v \n", block)
	}
}
