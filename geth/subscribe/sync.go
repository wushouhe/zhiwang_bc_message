package subscribe

import (
	"github.com/ethereum/go-ethereum/rpc"
	"context"
	"github.com/golang/glog"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"math/big"
	"zhiwang_bc_message/geth/blockdb"
	"database/sql"
)

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

func getLastBlock(client *rpc.Client) *big.Int {
	var block json.JsonHeader = json.JsonHeader{}
	if err := client.CallContext(context.Background(), &block, "eth_getBlockByNumber", toBlockNumArg(nil), true); err != nil {
		glog.Infof("call getBlockByNumber error: %v", err)
		return nil
	}

	return block.Number.ToInt()
}

/**
	同步所有数据
	问题：1，订阅之前获取最新区块，会漏取区块
	     2，订阅之后获取最新区块链 会多取区块（解决：插入时候去重）

	     方法已弃用
 */
func SyncAndSubscribBlock(client *rpc.Client, blockChan chan *json.JsonHeader) {

	//subscribe
	ListenNewBlock(client, blockChan)
	//订阅之后last blcok
	lastBlockNum := getLastBlock(client)
	//sync
	FillBlockRange(client, blockChan, int64(0), lastBlockNum.Int64())
}

func SyncBlocks(client *rpc.Client, db *sql.DB, blockChan chan *json.JsonHeader) {
	currentBlockNum := blockdb.LastBlockNumber(db)
	if currentBlockNum != 0 {
		currentBlockNum = currentBlockNum + 1
	}
	lastBlockNum := getLastBlock(client)

	glog.Infof("current %#v last %#v \n", currentBlockNum, lastBlockNum.Int64())
	//sync
	if currentBlockNum <= lastBlockNum.Int64() {
		FillBlockRange(client, blockChan, currentBlockNum, lastBlockNum.Int64())
		glog.Infof("同步完成 \n")
	}

}

/////////////////////////////////////
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
		batchRequest(client, blockChan, start, i)
		if i >= end {
			glog.Infof("i>=end i=%#v,end=%#v", i, end)
			break
		}
		start = i + 1
	}
}

func batchRequest(client *rpc.Client, blockChan chan *json.JsonHeader, start, end int64) {
	length := end - start + 1
	reqs := make([]rpc.BatchElem, length)
	//request
	for i := range reqs {
		/*reqs[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{toBlockNumArg(big.NewInt(start + int64(i))), true},
			Result: &json.JsonHeader{},
		}*/
		reqs[i].Method = "eth_getBlockByNumber"
		reqs[i].Args = []interface{}{toBlockNumArg(big.NewInt(start + int64(i))), true}
		reqs[i].Result = &json.JsonHeader{}
	}

	if err := client.BatchCall(reqs); err != nil {
		glog.Errorf("err %v", err)
	}

	//go func(rs []rpc.BatchElem) {
	//
	//	for _, req := range rs {
	//		blockChan <- req.Result.(*json.JsonHeader)
	//	}
	//}(reqs)
	for _, req := range reqs {
		blockChan <- req.Result.(*json.JsonHeader)
	}
	reqs = nil

}


//补全缺失区块
func FillBlockRangeComplete(client *rpc.Client, blockChan chan *json.JsonHeader, list []int64) {
	start := int64(0)
	end := int64(len(list) - 1)

	step := int64(99)
	i := int64(0)
	for {
		i = start + step
		if i > end {
			i = end
		}

		//loop
		batchRequestComplete(client, blockChan, start, i, list)
		if i >= end {
			glog.Infof("i>=end i=%#v,end=%#v", i, end)
			break
		}
		start = i + 1
	}
}

func batchRequestComplete(client *rpc.Client, blockChan chan *json.JsonHeader, start, end int64, list []int64) {
	length := end - start + 1
	reqs := make([]rpc.BatchElem, length)
	//request
	for i := range reqs {

		reqs[i].Method = "eth_getBlockByNumber"
		reqs[i].Args = []interface{}{toBlockNumArg(big.NewInt(list[start + int64(i)])), true}
		reqs[i].Result = &json.JsonHeader{}
	}

	if err := client.BatchCall(reqs); err != nil {
		glog.Errorf("err %v", err)
	}

	for _, req := range reqs {
		blockChan <- req.Result.(*json.JsonHeader)
	}
	reqs = nil

}