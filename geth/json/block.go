package json

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"fmt"
)

type JsonTransaction struct {
	BlockHash        *common.Hash                `json:blockHash`
	BlockNumber      *hexutil.Big        `json:blockNumber`
	From             *common.Address        `json:from`
	Hash             *common.Hash        `json:"hash"`
	TransactionIndex *hexutil.Uint64        `json:transactionIndex`
	AccountNonce     *hexutil.Uint64        `json:"nonce"`
	GasPrice         *hexutil.Big        `json:"gasPrice"`
	Gas              *hexutil.Big        `json:"gas"`
	Recipient        *common.Address        `json:"to"`
	Amount           *hexutil.Big        `json:"value"`
	Payload          *hexutil.Bytes        `json:"input"`
	V                *hexutil.Big        `json:"v"`
	R                *hexutil.Big        `json:"r"`
	S                *hexutil.Big        `json:"s"`
}

func (t *JsonTransaction) String() string {
	str := fmt.Sprintf(`
			BlockHash %s
			BlockNumber %s
			From %s
			Hash %s
			TransactionIndex %s
			AccountNonce %s
			GasPrice %s
			Gas %s
			Recipient %s
			Amount %s
			Payload %v
			V 0x%x
			R 0x%x
			S 0x%x `,
		t.BlockHash.Hex(),
		t.BlockNumber.ToInt(),
		t.From.Hex(),
		t.Hash.Hex(),
		t.TransactionIndex,
		t.AccountNonce,
		t.GasPrice,
		t.Gas,
		t.Recipient.Hex(),
		t.Amount,
		t.Payload,
		t.V,
		t.R,
		t.S)

	return str
}

type JsonHeader struct {
	ParentHash      *common.Hash        `json:"parentHash"`
	UncleHash       *common.Hash        `json:"sha3Uncles"`
	Coinbase        *common.Address        `json:"miner"`
	Root            *common.Hash        `json:"stateRoot"`
	TxHash          *common.Hash        `json:"transactionsRoot"`
	ReceiptHash     *common.Hash        `json:"receiptsRoot"`
	Bloom           *types.Bloom                `json:"logsBloom"`
	Difficulty      *hexutil.Big        `json:"difficulty"`
	Number          *hexutil.Big        `json:"number"`
	GasLimit        *hexutil.Big        `json:"gasLimit"`
	GasUsed         *hexutil.Big        `json:"gasUsed"`
	Time            *hexutil.Big        `json:"timestamp"`
	Extra           *hexutil.Bytes        `json:"extraData"`
	MixDigest       *common.Hash        `json:"mixHash"`
	Nonce           *types.BlockNonce        `json:"nonce"`
	Transactions    []*JsonTransaction        `json:transactions`
	Size            *hexutil.Big                `json:size`
	TotalDifficulty *hexutil.Big                `json:totalDifficulty`
	Hash            *common.Hash                `json:hash`
}

func (block *JsonHeader) String() string {
	str := fmt.Sprintf(`
			coinbase %s
			number %s
			difficulty %s
			gaslimit %s
			gasused %s
			nonce %#v
			parenthash %s
			txhash %s
			receipthash %s
			unclehash %s
			bloom %s
			extra %s
			MixDigest %s
			root %s
			time %s
			size %#x
			totalDifficultyv %#x
			hash %s
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
		block.UncleHash.Hex(),
		block.Bloom.Big(),
		block.Extra.String(),
		block.MixDigest.Hex(),
		block.Root.Hex(),
		block.Time.ToInt(),
		block.Size.ToInt(),
		block.TotalDifficulty.ToInt(),
		block.Hash.Hex(),
		block.Transactions)

	return str
}

var (
	errMissingHeaderMixDigest = errors.New("missing mixHash in JSON block header")
	errMissingHeaderFields = errors.New("missing required JSON block header fields")
)

func UnmarshalJSON(h *types.Header, dec *JsonHeader) error {

	// Ensure that all fields are set. MixDigest is checked separately because
	// it is a recent addition to the spec (as of August 2016) and older RPC server
	// implementations might not provide it.
	if dec.MixDigest == nil {
		return errMissingHeaderMixDigest
	}
	if dec.ParentHash == nil || dec.UncleHash == nil || dec.Coinbase == nil ||
		dec.Root == nil || dec.TxHash == nil || dec.ReceiptHash == nil ||
		dec.Bloom == nil || dec.Difficulty == nil || dec.Number == nil ||
		dec.GasLimit == nil || dec.GasUsed == nil || dec.Time == nil ||
		dec.Extra == nil || dec.Nonce == nil {
		return errMissingHeaderFields
	}
	// Assign all values.
	h.ParentHash = *dec.ParentHash
	h.UncleHash = *dec.UncleHash
	h.Coinbase = *dec.Coinbase
	h.Root = *dec.Root
	h.TxHash = *dec.TxHash
	h.ReceiptHash = *dec.ReceiptHash
	h.Bloom = *dec.Bloom
	h.Difficulty = (*big.Int)(dec.Difficulty)
	h.Number = (*big.Int)(dec.Number)
	h.GasLimit = (*big.Int)(dec.GasLimit)
	h.GasUsed = (*big.Int)(dec.GasUsed)
	h.Time = (*big.Int)(dec.Time)
	h.Extra = *dec.Extra
	h.MixDigest = *dec.MixDigest
	h.Nonce = *dec.Nonce
	return nil
}