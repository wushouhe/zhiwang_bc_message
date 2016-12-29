package main

import (
	"fmt"
	"time"
	"errors"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type jsonHeader struct {
	ParentHash  *common.Hash    `json:"parentHash"`
	UncleHash   *common.Hash    `json:"sha3Uncles"`
	Coinbase    *common.Address `json:"miner"`
	Root        *common.Hash    `json:"stateRoot"`
	TxHash      *common.Hash    `json:"transactionsRoot"`
	ReceiptHash *common.Hash    `json:"receiptsRoot"`
	Bloom       *types.Bloom          `json:"logsBloom"`
	Difficulty  *hexutil.Big    `json:"difficulty"`
	Number      *hexutil.Big    `json:"number"`
	GasLimit    *hexutil.Big    `json:"gasLimit"`
	GasUsed     *hexutil.Big    `json:"gasUsed"`
	Time        *hexutil.Big    `json:"timestamp"`
	Extra       *hexutil.Bytes  `json:"extraData"`
	MixDigest   *common.Hash    `json:"mixHash"`
	Nonce       *types.BlockNonce     `json:"nonce"`
}

var (
	errMissingHeaderMixDigest = errors.New("missing mixHash in JSON block header")
	errMissingHeaderFields = errors.New("missing required JSON block header fields")
)

func UnmarshalJSON(h *types.Header,dec *jsonHeader) error {

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

func main() {

	client, _ := rpc.Dial("http://172.16.10.163:8545")

	filterId, _ := subscribeNewBlock(client)
	defer unSubscribeNewBlock(client, filterId)

	blockIdChan := make(chan []string)
	go getNewBlockIds(client, filterId, blockIdChan)

	blockChan := make(chan *jsonHeader)
	go blockDetail(client, blockIdChan, blockChan)

	for {
		select {
		case block := <-blockChan:
			var h *types.Header=&types.Header{}
			UnmarshalJSON(h,block)
			fmt.Println(h.Coinbase)
			fmt.Println(h.Nonce)
		}
	}

}

func getNewBlockIds(client *rpc.Client, filterId string, blockIdChan chan []string) {
	timer := time.NewTimer(time.Second * 10)
	var result []string = make([]string, 0)
	for {
		select {
		case <-timer.C:
			if err := client.CallContext(context.Background(), &result, "eth_getFilterChanges", filterId); err != nil {
				fmt.Println("can't get latest block:", err)
				return
			}
			blockIdChan <- result
			timer.Reset(time.Second * 10)
		}
	}
}

func subscribeNewBlock(client *rpc.Client) (string, error) {
	var result string
	if err := client.CallContext(context.Background(), &result, "eth_newBlockFilter"); err != nil {
		fmt.Println("can't get latest block:", err)
		return "", err
	}
	return result, nil
}

func unSubscribeNewBlock(client *rpc.Client, filterId string) (bool, error) {
	var result bool
	if err := client.CallContext(context.Background(), &result, "eth_uninstallFilter"); err != nil {
		fmt.Println("can't get latest block:", err)
		return false, err
	}
	return result, nil
}

func blockDetail(client *rpc.Client, blockIdChan chan []string, blockChan chan *jsonHeader) {
	for {
		select {
		case blockIds := <-blockIdChan:
			for _, blockId := range blockIds {
				block, err := blockDetailByHash(client, blockId)
				if err != nil {
					fmt.Println(err)
					return
				}
				blockChan <- block
			}
		}
	}
}
func blockDetailByHash(client *rpc.Client, blockId string) (*jsonHeader, error) {
	var block jsonHeader = jsonHeader{}
	if err := client.CallContext(context.Background(), &block, "eth_getBlockByHash", blockId, true); err != nil {
		fmt.Println("can't get latest block:", err)
		return nil, err
	}
	return &block, nil
}



