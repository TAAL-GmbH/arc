package store

import (
	"github.com/bitcoin-sv/arc/internal/blocktx/blocktx_api"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
)

type BlockGap struct {
	Height uint64
	Hash   *chainhash.Hash
}

type TxWithMerklePath struct {
	Hash       []byte
	MerklePath string
}

type GetMinedTransactionResult struct {
	TxHash      []byte
	BlockHash   []byte
	BlockHeight uint64
	MerklePath  string
}

type BlockStatusUpdate struct {
	Hash   []byte
	Status blocktx_api.Status
}
