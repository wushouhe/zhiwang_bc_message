package subscribe

import (
	"github.com/ethereum/go-ethereum/rpc"
	"context"
	"github.com/ethereum/go-ethereum/logger/glog"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"math/big"
	//"zhiwang_bc_message/geth/utils"
)
//同步区块 不再使用，会造成以下错误：
//call getBlockByNumber error: Post http://172.16.10.163:8545: dial tcp 172.16.10.163:8545: bind: An operation on a socket could not be performed because the system lacked sufficient buffer space or because a queue was full.nil block
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
		fmt.Printf("call getBlockByNumber error: %v", err)
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
	FillBlockRange(client, blockChan, int64(0), lastBlockNum.Int64())
}

/////////////////////////////////////
/**
request 要复用 与分配99个

 */

func BatchRequest(client *rpc.Client, blockChan chan *json.JsonHeader, start, end int64) {
	length := end - start + 1
	reqs := make([]rpc.BatchElem, length)


	//request
	for i := range reqs {
		reqs[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{toBlockNumArg(big.NewInt(start + int64(i))), true},
			Result: &json.JsonHeader{},
		}
	}

	if err := client.BatchCall(reqs); err != nil {
		fmt.Printf("err %v", err)
	}

	go func(rs []rpc.BatchElem) {

		for _, req := range rs {
			blockChan <- req.Result.(*json.JsonHeader)
			//utils.PrintBlock(req.Result.(*json.JsonHeader))
		}
	}(reqs)

}

//批量获取[start,end]之间的区块
func FillBlockRange(client *rpc.Client, blockChan chan *json.JsonHeader, start, end int64) {
	step := int64(99)
	i := int64(0)
	for {
		i = start + step
		if i > end {
			i = end
		}

		//loop
		BatchRequest(client, blockChan, start, i)
		if i >= end {
			break
		}
		start = i + 1
	}

}