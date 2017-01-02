package main

import (
	"zhiwang_bc_message/geth/blockdb"
	"fmt"
)

func main() {
	db:=blockdb.NewDB()
	n:=blockdb.LastBlockNumber(db)
	fmt.Println(n)
}
