package blockdb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"zhiwang_bc_message/geth/json"
)

var (
	ip = "172.16.10.162:3306"//IP地址
	username = "root"//用户名
	passwd = "123456"//密码
	dbname = "zw_bc"//库名
)

func NewDB() *sql.DB{
	db, _ := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, passwd, ip, dbname))
	return db
}


func InesrtBlockChan(db *sql.DB, block *json.JsonHeader) {
	InsertBlock(db, block)

	InsertTransactions(db,block.Transactions)
	/*for _, tx := range block.Transactions {
		InsertTransaction(db, tx)
	}*/

}

