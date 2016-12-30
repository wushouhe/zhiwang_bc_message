package utils

import (
	"fmt"
	"zhiwang_bc_message/geth/json"
)

func PrintBlock(block *json.JsonHeader) {
	fmt.Printf(`
			coinbase %s
			number %s
			difficulty %s
			gaslimit %s
			gasused %s
			nonce %#v
			parenthash %s
			txhash %s
			receipthash %s
			bloom %s
			extra %s
			MixDigest %s
			root %s
			time %s
			Transactions %v `,
		block.Coinbase.Hex(),
		block.Number.ToInt(),
		block.Difficulty.ToInt(),
		block.GasLimit.ToInt(),
		block.GasUsed.ToInt(),
		block.Nonce.Uint64(),
		block.ParentHash.Hex(),
		block.TxHash.Hex(),
		block.ReceiptHash.Hex(),
		block.Bloom.Big(),
		block.Extra.String(),
		block.MixDigest.Hex(),
		block.Root.Hex(),
		block.Time.ToInt(),
		block.Transactions)
}

