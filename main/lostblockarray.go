package main

import (
	"zhiwang_bc_message/geth/lostblock"
	"flag"
)

func main() {
	flag.Parse()
	lostblock.TestBinary()
}
