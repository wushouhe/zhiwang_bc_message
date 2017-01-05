package blockdb

import (
	"database/sql"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"strings"
	"strconv"
	"github.com/golang/glog"
)

func InsertBlock(db *sql.DB, block *json.JsonHeader) {

	stmt, e := db.Prepare(`INSERT INTO blocks (hash, parentHash, nonce, number, extraData, gasLimit, gasUsed, miner, mixHash, receiptsRoot,
	 stateRoot, sha3Uncles, logsBloom, size, difficulty, totalDifficulty, timestamp, transactionsRoot)
	 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if e != nil {
		glog.Infof("stmt %v \n", e)
	}
	rows, err := stmt.Query(transBlock(block)...)
	defer stmt.Close()
	if err != nil {
		glog.Infof("insert data error: %v\n", err)
	}
	rows.Close()
}

func InsertBlockBatch(tx *sql.Tx, block *json.JsonHeader) {

	stmt, e := tx.Prepare(`INSERT INTO blocks (hash, parentHash, nonce, number, extraData, gasLimit, gasUsed, miner, mixHash, receiptsRoot,
	 stateRoot, sha3Uncles, logsBloom, size, difficulty, totalDifficulty, timestamp, transactionsRoot)
	 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if e != nil {
		glog.Infof("stmt %v \n", e)
	}
	_, err := stmt.Exec(transBlock(block)...)
	defer stmt.Close()
	if err != nil {
		glog.Infof("insert data error: %v\n", err)
	}
}

func RemoveRepeatRows(db *sql.DB) {
	stmt, err := db.Prepare(`delete from blocks
				where number in (
						select a.number
						from
						(
							select	number
							from blocks
							group by number
							having count(number) > 1
						) a
					)`)
	defer stmt.Close()
	if err != nil {
		glog.Infof("stmt %v \n", err)
	}
	rows, err := stmt.Exec()
	if err != nil {
		glog.Infof("insert data error: %v\n", err)
	}
	affectedCount, _ := rows.RowsAffected()
	glog.Infof("affect rows %d \n", affectedCount)
}

func LastBlockNumber(db *sql.DB) int64 {
	var lastNumber string
	row := db.QueryRow("select max(number)  from blocks ")
	err := row.Scan(&lastNumber)
	if err != nil {
		glog.Infof("last number err %v \n", err)
		lastNumber = "0"
	}
	i, err := strconv.ParseInt(lastNumber, 10, 64)
	if err != nil {
		glog.Infof("parse int64 error %v \n", err)
		return int64(0)
	}
	glog.Infof("last number from db is %s \n", lastNumber)
	return i
}


func MinBlockNumber(db *sql.DB) int64 {
	var lastNumber string
	row := db.QueryRow("select min(number)  from blocks ")
	err := row.Scan(&lastNumber)
	if err != nil {
		glog.Infof("min number err %v \n", err)
		lastNumber = "0"
	}
	i, err := strconv.ParseInt(lastNumber, 10, 64)
	if err != nil {
		glog.Infof("parse int64 error %v \n", err)
		return int64(0)
	}
	glog.Infof("min number from db is %s \n", lastNumber)
	return i
}


func transBlock(block *json.JsonHeader) []interface{} {
	blockStrs := strings.Split(blockStr(block), ",")
	blockObjs := make([]interface{}, len(blockStrs))
	for i, blockStr := range blockStrs {
		blockObjs[i] = blockStr
	}
	return blockObjs
}

func blockStr(block *json.JsonHeader) string {
	str := fmt.Sprintf(`%s,%s,%#v,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%#x,%s,%#x,%s,%s`, block.Hash.Hex(), block.ParentHash.Hex(), block.Nonce.Uint64(), block.Number.ToInt(), block.Extra.String(), block.GasLimit.ToInt(), block.GasUsed.ToInt(), block.Coinbase.Hex(), block.MixDigest.Hex(), block.ReceiptHash.Hex(), block.Root.Hex(), block.UncleHash.Hex(), block.Bloom.Big(), block.Size.ToInt(), block.Difficulty.ToInt(), block.TotalDifficulty.ToInt(), block.Time.ToInt(), block.TxHash.Hex())
	return str
}

