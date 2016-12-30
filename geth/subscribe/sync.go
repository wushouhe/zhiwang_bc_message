package subscribe

import (
	"github.com/ethereum/go-ethereum/rpc"
	"context"
	"github.com/ethereum/go-ethereum/logger/glog"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"math/big"
)

func SyncBlock(client *rpc.Client, from, to int64) []*json.JsonHeader {
	blocks := make([]*json.JsonHeader, 0)
	for i := from; i < to; i++ {
		blocks = append(blocks, getBlockByNumber(client, big.NewInt(i)))
	}
	return blocks
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
func lastBlockNum() {

}