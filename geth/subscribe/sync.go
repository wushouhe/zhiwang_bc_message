package subscribe

import (
	"github.com/ethereum/go-ethereum/rpc"
	"context"
	"github.com/ethereum/go-ethereum/logger/glog"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"math/big"
)

func SyncBlock(client *rpc.Client, blockChan chan *json.JsonHeader, from, to int64) {
	for i := from; i < to; i++ {
		go func(num int64) {
			blockChan <- getBlockByNumber(client, big.NewInt(num))
		}(i)
	}
}

func getBlockByNumber(client *rpc.Client, blockNumber *big.Int) (*json.JsonHeader) {
	var block json.JsonHeader = json.JsonHeader{}
	if err := client.CallContext(context.Background(), &block, "eth_getBlockByNumber", toBlockNumArg(blockNumber), true); err != nil {
		glog.Infof("call getBlockByNumber error: %v", err)
		return nil
	}
	return &block
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return fmt.Sprintf("%#x", number)
}

//chan
func getLastBlock(client *rpc.Client, blockChan chan *json.JsonHeader) *big.Int {
	var block json.JsonHeader = json.JsonHeader{}
	if err := client.CallContext(context.Background(), &block, "eth_getBlockByNumber", toBlockNumArg(nil), true); err != nil {
		glog.Infof("call getBlockByNumber error: %v", err)
		return nil
	}

	go func() {
		blockChan <- &block
	}()

	return block.Number.ToInt()
}

func checkSyncBlock() {

}