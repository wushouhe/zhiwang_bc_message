package blockdb

import (
	"database/sql"
	"log"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"strings"
)

func InsertTransaction(db *sql.DB, tx *json.JsonTransaction) {

	stmt, e := db.Prepare(`
	INSERT INTO transactions (hash,
			blockHash,
			blockNumber,
			tx_from,
			tx_to,
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
		?)`)
	if e != nil {
		fmt.Printf("tx err %v \n", e)
	}
	rows, err := stmt.Query(transTx(tx)...)
	defer stmt.Close()
	if err != nil {
		log.Fatalf("insert data error: %v\n", err)
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
	str := fmt.Sprintf(`%s,%s,%s,%s,%s,%s,%v,%s,%s,%s,%s,0x%x,0x%x,0x%x`, t.Hash.Hex(), t.BlockHash.Hex(), t.BlockNumber.ToInt(), t.From.Hex(), t.Recipient.Hex(), t.Amount, t.Payload, t.AccountNonce, t.TransactionIndex, t.Gas, t.GasPrice, t.V, t.R, t.S)
	return str
}