package store

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/arc/internal/blocktx/blocktx_api"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
)

var (
	ErrNotFound                         = errors.New("not found")
	ErrBlockProcessingDuplicateKey      = errors.New("block hash already exists")
	ErrBlockNotFound                    = errors.New("block not found")
	ErrUnableToPrepareStatement         = errors.New("unable to prepare statement")
	ErrUnableToDeleteRows               = errors.New("unable to delete rows")
	ErrFailedToInsertBlock              = errors.New("failed to insert block")
	ErrFailedToOpenDB                   = errors.New("failed to open postgres database")
	ErrFailedToInsertTransactions       = errors.New("failed to bulk insert transactions")
	ErrFailedToGetRows                  = errors.New("failed to get rows")
	ErrFailedToSetBlockProcessing       = errors.New("failed to set block processing")
	ErrFailedToExecuteTxUpdateQuery     = errors.New("failed to execute transaction update query")
	ErrMismatchedTxsAndMerklePathLength = errors.New("mismatched transactions and merkle path length")
)

type BlocktxStore interface {
	RegisterTransactions(ctx context.Context, txHashes [][]byte) (updatedTxs []*chainhash.Hash, err error)
	GetBlock(ctx context.Context, hash *chainhash.Hash) (*blocktx_api.Block, error)
	GetBlockByHeight(ctx context.Context, height uint64, status blocktx_api.Status) (*blocktx_api.Block, error)
	GetChainTip(ctx context.Context) (*blocktx_api.Block, error)
	UpsertBlock(ctx context.Context, block *blocktx_api.Block) (uint64, error)
	UpsertBlockTransactions(ctx context.Context, blockID uint64, txsWithMerklePaths []TxWithMerklePath) (registeredTxs []TxWithMerklePath, err error)
	MarkBlockAsDone(ctx context.Context, hash *chainhash.Hash, size uint64, txCount uint64) error
	GetBlockGaps(ctx context.Context, heightRange int) ([]*BlockGap, error)
	ClearBlocktxTable(ctx context.Context, retentionDays int32, table string) (*blocktx_api.RowsAffectedResponse, error)
	GetMinedTransactions(ctx context.Context, hashes [][]byte, blockStatus blocktx_api.Status) ([]GetMinedTransactionResult, error)
	GetLongestChainFromHeight(ctx context.Context, height uint64) ([]*blocktx_api.Block, error)
	GetStaleChainBackFromHash(ctx context.Context, hash []byte) ([]*blocktx_api.Block, error)
	UpdateBlocksStatuses(ctx context.Context, blockStatusUpdates []BlockStatusUpdate) error

	SetBlockProcessing(ctx context.Context, hash *chainhash.Hash, processedBy string) (string, error)
	GetBlockHashesProcessingInProgress(ctx context.Context, processedBy string) ([]*chainhash.Hash, error)
	DelBlockProcessing(ctx context.Context, hash *chainhash.Hash, processedBy string) (int64, error)
	VerifyMerkleRoots(ctx context.Context, merkleRoots []*blocktx_api.MerkleRootVerificationRequest, maxAllowedBlockHeightMismatch int) (*blocktx_api.MerkleRootVerificationResponse, error)

	Ping(ctx context.Context) error
	Close() error
}
