package json

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type JsonHeader struct {
	ParentHash   *common.Hash    `json:"parentHash"`
	UncleHash    *common.Hash    `json:"sha3Uncles"`
	Coinbase     *common.Address `json:"miner"`
	Root         *common.Hash    `json:"stateRoot"`
	TxHash       *common.Hash    `json:"transactionsRoot"`
	ReceiptHash  *common.Hash    `json:"receiptsRoot"`
	Bloom        *types.Bloom          `json:"logsBloom"`
	Difficulty   *hexutil.Big    `json:"difficulty"`
	Number       *hexutil.Big    `json:"number"`
	GasLimit     *hexutil.Big    `json:"gasLimit"`
	GasUsed      *hexutil.Big    `json:"gasUsed"`
	Time         *hexutil.Big    `json:"timestamp"`
	Extra        *hexutil.Bytes  `json:"extraData"`
	MixDigest    *common.Hash    `json:"mixHash"`
	Nonce        *types.BlockNonce     `json:"nonce"`
	Transactions types.Transactions    `json:transactions`
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