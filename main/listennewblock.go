package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/blockdb"
	"github.com/golang/glog"
	"sync"
	"flag"
	"zhiwang_bc_message/geth/lostblock"
)

var once sync.Once

func main() {
	flag.Parse()
	//172.16.10.163
	//139.196.178.168
	client, _ := rpc.Dial("http://139.196.178.168:8545")
	blockChan := make(chan *json.JsonHeader, 1000)
	glog.Infof("正在监听new Block.........")
	subscribe.ListenNewBlock(client, blockChan)
	db := blockdb.NewDB()
	lostBlockFunc := func() {
		lostblock.SyncLostBlock(client,db,blockChan)
	}
	/*for {
		block := <-blockChan
		blockdb.InesrtBlockChan(db, block)
		glog.Infof("block %v ", block)
		once.Do(lostBlockFunc)
	}*/

	var w sync.WaitGroup
	w.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			for {
				block := <-blockChan
				glog.Infof("importing block %s \n", block.Number.ToInt())
				blockdb.InesrtBlockChan(db, block)
				once.Do(lostBlockFunc)
			}
			w.Done()
		}()
	}
	w.Wait()

}





