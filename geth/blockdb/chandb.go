package blockdb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"zhiwang_bc_message/geth/json"
	"github.com/golang/glog"
	. "zhiwang_bc_message/geth/config"
)



func NewDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", Cfg.Mysql.Username, Cfg.Mysql.Passwd, Cfg.Mysql.Ip, Cfg.Mysql.Port, Cfg.Mysql.BaseName))
	if err != nil {
		glog.Infof("create db error %v ", err)
	}
	return db
}

func InesrtBlockChan(db *sql.DB, block *json.JsonHeader) {
	tx, err := db.Begin()
	if err != nil {
		glog.Errorf("tx error %v \n", err)
	}
	InsertBlockBatch(tx, block)

	InsertTransactionsBatch(tx, block.Transactions)
	tx.Commit()
	/*for _, tx := range block.Transactions {
		InsertTransaction(db, tx)
	}*/

}

