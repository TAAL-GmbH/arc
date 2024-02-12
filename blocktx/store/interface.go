package store

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/arc/blocktx/blocktx_api"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
)

var (
	ErrBlockNotFound = errors.New("block not found")
)

type BlocktxStore interface {
	RegisterTransactions(ctx context.Context, transaction []*blocktx_api.TransactionAndSource) error
	TryToBecomePrimary(ctx context.Context, myHostName string) error
	GetPrimary(ctx context.Context) (string, error)
	GetBlock(ctx context.Context, hash *chainhash.Hash) (*blocktx_api.Block, error)
	InsertBlock(ctx context.Context, block *blocktx_api.Block) (uint64, error)
	UpdateBlockTransactions(ctx context.Context, blockId uint64, transactions []*blocktx_api.TransactionAndSource, merklePaths []string) ([]UpdateBlockTransactionsResult, error)
	MarkBlockAsDone(ctx context.Context, hash *chainhash.Hash, size uint64, txCount uint64) error
	GetBlockGaps(ctx context.Context, heightRange int) ([]*BlockGap, error)
	ClearBlocktxTable(ctx context.Context, retentionDays int32, table string) (*blocktx_api.ClearDataResponse, error)
	Close() error
}
