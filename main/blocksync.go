package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/subscribe"
	"fmt"
)

func main() {
	client, _ := rpc.Dial("http://172.16.10.163:8545")
	blocks := subscribe.SyncBlock(client, 1, 10)
	for k, v := range blocks {
		fmt.Printf("%d, %v \n", k, v)
	}
}
