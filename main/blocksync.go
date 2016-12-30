package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/json"
	"fmt"
)

func main() {
	client, _ := rpc.Dial("http://172.16.10.163:8545")
	blockChan := make(chan *json.JsonHeader)
	subscribe.SyncBlock(client, blockChan, 1, 10)
	for v := range blockChan {
		fmt.Printf(" %v \n", v)
	}
}
