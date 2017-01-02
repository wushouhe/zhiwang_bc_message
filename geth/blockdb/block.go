package blockdb

import (
	"database/sql"
	"log"
	"zhiwang_bc_message/geth/json"
	"fmt"
	"strings"
	"strconv"
)

func InsertBlock(db *sql.DB, block *json.JsonHeader) {

	stmt, e := db.Prepare(`INSERT INTO blocks (hash, parentHash, nonce, number, extraData, gasLimit, gasUsed, miner, mixHash, receiptsRoot,
	 stateRoot, sha3Uncles, logsBloom, size, difficulty, totalDifficulty, timestamp, transactionsRoot)
	 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if e != nil {
		fmt.Printf("stmt %v \n", e)
	}
	rows, err := stmt.Query(transBlock(block)...)
	defer stmt.Close()
	if err != nil {
		log.Fatalf("insert data error: %v\n", err)
	}
	rows.Close()
}

func LastBlockNumber(db *sql.DB) int64 {
	var lastNumber string
	row := db.QueryRow("select number  from blocks order by number desc limit 0,1")
	err := row.Scan(&lastNumber)
	if err != nil {
		fmt.Printf("last number err %v \n", err)
		lastNumber = "0"
	}
	i, err := strconv.ParseInt(lastNumber, 10, 64)
	if err != nil {
		fmt.Printf("parse int64 error %v \n", err)
		return int64(0)
	}
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
	fmt.Println(str)
	return str
}

