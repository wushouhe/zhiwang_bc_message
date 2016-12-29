package rpc

import (
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"context"
	"zhiwang_bc_message/geth/json"
	"github.com/ethereum/go-ethereum/logger/glog"
)

//入口
func ListenNewBlock(client *rpc.Client, blockChan chan *json.JsonHeader) {
	//订阅
	filterId, _ := SubscribeNewBlock(client)
	defer UnSubscribeNewBlock(client, filterId)

	//监听new block
	blockIdChan := make(chan []string)
	go GetNewBlockIds(client, filterId, blockIdChan)

	//block 信息
	go BlockDetail(client, blockIdChan, blockChan)

}

/**
根据订阅Id获取新生成block Id
filterId 订阅ID
blockIdChan new blockID 队列
 */
func GetNewBlockIds(client *rpc.Client, filterId string, blockIdChan chan []string) {
	timer := time.NewTimer(time.Second * 10)
	var result []string = make([]string, 0)
	for {
		select {
		case <-timer.C:
			if err := client.CallContext(context.Background(), &result, "eth_getFilterChanges", filterId); err != nil {
				glog.Infof("call eth_getFilterChanges error: %v", err)
				return
			}
			blockIdChan <- result
			timer.Reset(time.Second * 10)
		}
	}
}

/**
订阅事件（newBlock）
 */
func SubscribeNewBlock(client *rpc.Client) (string, error) {
	var result string
	if err := client.CallContext(context.Background(), &result, "eth_newBlockFilter"); err != nil {
		glog.Infof("call eth_newBlockFilter error: %v", err)
		return "", err
	}
	return result, nil
}

/**
取消订阅事件（newBlock）
filterId RPC ID
 */
func UnSubscribeNewBlock(client *rpc.Client, filterId string) (bool, error) {
	var result bool
	if err := client.CallContext(context.Background(), &result, "eth_uninstallFilter"); err != nil {
		glog.Infof("call eth_uninstallFilter error: %v", err)
		return false, err
	}
	return result, nil
}

/**
获取区块信息放入 blockchan中
blockIdChan  blockId 队列
blockChan 区块信息队列
 */
func BlockDetail(client *rpc.Client, blockIdChan chan []string, blockChan chan *json.JsonHeader) {
	for {
		select {
		case blockIds := <-blockIdChan:
			for _, blockId := range blockIds {
				block, err := BlockDetailByHash(client, blockId)
				if err != nil {
					glog.Infof("BlockDetail error: %v", err)
					return
				}
				blockChan <- block
			}
		}
	}
}

/**
通过 blockId(block hash)获取block 详细信息
blockId block hash
 */
func BlockDetailByHash(client *rpc.Client, blockId string) (*json.JsonHeader, error) {
	var block json.JsonHeader = json.JsonHeader{}
	if err := client.CallContext(context.Background(), &block, "eth_getBlockByHash", blockId, true); err != nil {
		glog.Infof("call eth_getBlockByHash error: %v", err)
		return nil, err
	}
	return &block, nil
}
