package lostblock

import (
	"zhiwang_bc_message/geth/blockdb"
	"database/sql"
	"github.com/golang/glog"
	"zhiwang_bc_message/geth/subscribe"
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
)

func SyncLostBlock(client *rpc.Client, db *sql.DB, blockChan chan *json.JsonHeader) {
	glog.Infof("删除重复区块")
	blockdb.RemoveRepeatRows(db)
	glog.Infof("检查缺失区块")
	lostList := MysqlLoop(db)
	go subscribe.FillBlockRangeComplete(client, blockChan, lostList)

	//最小区块号是否为零
	minBlockNumber := blockdb.MinBlockNumber(db)
	if minBlockNumber > int64(0) {
		glog.Infof("sync block [0,%d] ", minBlockNumber - 1)
		go subscribe.FillBlockRange(client, blockChan, int64(0), minBlockNumber - 1)
	}
}