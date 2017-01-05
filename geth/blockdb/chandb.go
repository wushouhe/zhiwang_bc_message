package blockdb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"zhiwang_bc_message/geth/json"
	"github.com/golang/glog"
)

var (
	/*ip = "172.16.10.162:3306"//IP地址
	username = "root"//用户名
	passwd = "123456"//密码
	dbname = "zw_bc"//库名*/

	ip = "127.0.0.1:3306"//IP地址
	username = "root"//用户名
	passwd = "root"//密码
	dbname = "ethereum"//库名
)

func NewDB() *sql.DB {
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, passwd, ip, dbname))
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

