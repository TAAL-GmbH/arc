package store

import (
	"github.com/libsv/go-p2p/chaincfg/chainhash"
)

type BlockGap struct {
	Height uint64
	Hash   *chainhash.Hash
}

type UpdateBlockTransactionsResult struct {
	TxHash     []byte
	MerklePath string
}

type GetMinedTransactionResult struct {
	BlockHash   []byte
	BlockHeight uint64
	MerklePath  string
}
