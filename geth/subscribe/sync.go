package subscribe

import (
	"github.com/ethereum/go-ethereum/rpc"
	"context"
	"github.com/ethereum/go-ethereum/logger/glog"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"math/big"
)
//同步区块
func SyncBlock(client *rpc.Client, blockChan chan *json.JsonHeader, from, to int64) {
	for i := from; i < to; i++ {
		go func(num int64) {
			Loop:
			for {
				select {
				case blockChan <- getBlockByNumber(client, big.NewInt(num)):
					break Loop
				}
			}

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

func getLastBlock(client *rpc.Client, blockChan chan *json.JsonHeader) *big.Int {
	var block json.JsonHeader = json.JsonHeader{}
	if err := client.CallContext(context.Background(), &block, "eth_getBlockByNumber", toBlockNumArg(nil), true); err != nil {
		glog.Infof("call getBlockByNumber error: %v", err)
		return nil
	}

	go func() {
		Loop:
		for {
			select {
			case blockChan <- &block:
				break Loop
			}
		}

	}()

	return block.Number.ToInt()
}

/**
	同步所有数据
	问题：1，订阅之前获取最新区块，会漏取区块
	     2，订阅之后获取最新区块链 会多取区块（解决：插入时候去重）
 */
func SyncAndSubscribBlock(client *rpc.Client, blockChan chan *json.JsonHeader) {

	//subscribe
	ListenNewBlock(client, blockChan)
	//订阅之后last blcok
	lastBlockNum := getLastBlock(client, blockChan)
	//sync
	SyncBlock(client, blockChan, 0, lastBlockNum.Int64())
}
