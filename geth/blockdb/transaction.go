package blockdb

import (
	"database/sql"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"strings"
	"github.com/golang/glog"
	"github.com/ethereum/go-ethereum/common"
)

func InsertTransactions(db *sql.DB, txs []*json.JsonTransaction) {
	tx, err := db.Begin()
	if err != nil {
		glog.Errorf("tx error %v \n", err)
	}
	stmt, err := tx.Prepare(`
	INSERT INTO transactions (hash,
			blockHash,
			blockNumber,
			tx_from,
			tx_to,
			isContract,
			value,
			input,
			nonce,
			transactionIndex,
			gas,
			gasPrice,
			v,
			r,
			s)
	 VALUES (?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?)`)
	if err != nil {
		glog.Errorf("stmt err %v \n", err)
	}
	for _, tx := range txs {
		stmt.Exec(transTx(tx)...)
	}
	tx.Commit()
	defer stmt.Close()

}

func InsertTransactionsBatch(tx *sql.Tx, txs []*json.JsonTransaction) {
	stmt, err := tx.Prepare(`
	INSERT INTO transactions (hash,
			blockHash,
			blockNumber,
			tx_from,
			tx_to,
			isContract,
			value,
			input,
			nonce,
			transactionIndex,
			gas,
			gasPrice,
			v,
			r,
			s)
	 VALUES (?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?)`)
	if err != nil {
		glog.Errorf("stmt err %v \n", err)
	}
	for _, tx := range txs {
		stmt.Exec(transTx(tx)...)
	}
	defer stmt.Close()

}

func InsertTransaction(db *sql.DB, tx *json.JsonTransaction) {

	stmt, e := db.Prepare(`
	INSERT INTO transactions (hash,
			blockHash,
			blockNumber,
			tx_from,
			tx_to,
			isContract,
			value,
			input,
			nonce,
			transactionIndex,
			gas,
			gasPrice,
			v,
			r,
			s)
	 VALUES (?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?)`)
	if e != nil {
		glog.Errorf("tx err %v \n", e)
	}

	rows, err := stmt.Query(transTx(tx)...)
	defer stmt.Close()
	if err != nil {
		glog.Errorf("insert data error: %v\n", err)
	}
	rows.Close()

}

func transTx(tx *json.JsonTransaction) []interface{} {
	txStrs := strings.Split(txStr(tx), ",")
	txObjs := make([]interface{}, len(txStrs))
	for i, txStr := range txStrs {
		txObjs[i] = txStr
	}
	return txObjs
}

func txStr(t *json.JsonTransaction) string {
	var str string
	var isContract string = "false"
	if t.Recipient == nil {
		t.Recipient = &common.Address{}
		isContract = "true"
	}
	str = fmt.Sprintf(`%s,%s,%s,%s,%s,%s,%s,%v,%s,%s,%s,%s,0x%x,0x%x,0x%x`, t.Hash.Hex(), t.BlockHash.Hex(), t.BlockNumber.ToInt(), t.From.Hex(), t.Recipient.Hex(),isContract, t.Amount, t.Payload, t.AccountNonce, t.TransactionIndex, t.Gas, t.GasPrice, t.V, t.R, t.S)
	return str
}