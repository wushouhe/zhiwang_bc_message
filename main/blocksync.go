package main

import (
	"github.com/ethereum/go-ethereum/rpc"
	"zhiwang_bc_message/geth/subscribe"
	"zhiwang_bc_message/geth/json"
	"zhiwang_bc_message/geth/blockdb"
	"time"
	"github.com/golang/glog"
	"flag"
	"zhiwang_bc_message/geth/lostblock"
	"sync"
)

func main() {
	flag.Parse()
	client, _ := rpc.Dial("http://139.196.178.168:8545")
	blockChan := make(chan *json.JsonHeader, 1000)
	db := blockdb.NewDB()
	lostblock.SyncLostBlock(client, db, blockChan)
	glog.Info("开始同步区块...")
	go func() {
		for {
			subscribe.SyncBlocks(client, db, blockChan)
			time.Sleep(2 * time.Minute)
		}
	}()

	/*timer := time.NewTimer(time.Second * 120)
	for {
		select {
		case <-timer.C:
			subscribe.SyncBlocks(client, db, blockChan)
			timer.Reset(time.Second * 120)
		case block := <-blockChan:
			glog.Infof("importing block %s \n", block.Number.ToInt())
			blockdb.InesrtBlockChan(db, block)
			timer.Reset(time.Second * 120)
		}
	}*/

	var w sync.WaitGroup
	w.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			for {
				block := <-blockChan
				glog.Infof("importing block %s \n", block.Number.ToInt())
				blockdb.InesrtBlockChan(db, block)
			}
			w.Done()
		}()
	}
	w.Wait()


	/*for block := range blockChan {
		//glog.Infof("%v \n", block)
		glog.Infof("before insert into blockchan %v \n",time.Now())
		blockdb.InesrtBlockChan(db, block)
		glog.Infof("after insert into blockchan %v \n",time.Now())
	}*/

}
