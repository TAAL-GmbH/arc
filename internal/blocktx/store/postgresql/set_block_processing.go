package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"github.com/bitcoin-sv/arc/internal/tracing"

	"github.com/lib/pq"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
)

// SetBlockProcessing tries to insert a record to the block processing table in order to mark a certain block as being processed by a blocktx instance. If there is an entry from any block tx instance not older than `lockTime` then a new entry will be inserted successfully returning the name of the blocktx instance which has called the function
func (p *PostgreSQL) SetBlockProcessing(ctx context.Context, hash *chainhash.Hash, setProcessedBy string, lockTime time.Duration) (string, error) {
	// Try to set a block as being processed by this instance
	qInsert := `
		INSERT INTO blocktx.block_processing (block_hash, processed_by)
		SELECT $1, $2
		WHERE NOT EXISTS (
		  SELECT 1 FROM blocktx.block_processing bp WHERE bp.block_hash = $1 AND inserted_at > $3
		)
		RETURNING processed_by
	`

	// Todo: insert only if not setProcessedBy instance is not already processing X blocks

	var processedBy string
	err := p.db.QueryRowContext(ctx, qInsert, hash[:], setProcessedBy, p.now().Add(-1*lockTime)).Scan(&processedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var currentlyProcessedBy string
			err = p.db.QueryRowContext(ctx, `SELECT processed_by FROM blocktx.block_processing WHERE block_hash = $1 ORDER BY inserted_at DESC LIMIT 1`, hash[:]).Scan(&currentlyProcessedBy)
			if err != nil {
				return "", errors.Join(store.ErrBlockProcessingDuplicateKey, err)
			}
			return currentlyProcessedBy, store.ErrBlockProcessingDuplicateKey
		}

		return "", errors.Join(store.ErrFailedToSetBlockProcessing, err)
	}

	return processedBy, nil
}

func (p *PostgreSQL) DelBlockProcessing(ctx context.Context, hash *chainhash.Hash, processedBy string) (rowsAffected int64, err error) {
	ctx, span := tracing.StartTracing(ctx, "DelBlockProcessing", p.tracingEnabled, p.tracingAttributes...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	q := `
		DELETE FROM blocktx.block_processing WHERE block_hash = $1 AND processed_by = $2;
  `

	res, err := p.db.ExecContext(ctx, q, hash[:], processedBy)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ = res.RowsAffected()
	if rowsAffected != 1 {
		return 0, store.ErrBlockNotFound
	}

	return rowsAffected, nil
}

func (p *PostgreSQL) GetBlockHashesProcessingInProgress(ctx context.Context, processedBy string) ([]*chainhash.Hash, error) {
	// Check how many blocks this instance is currently processing
	q := `
		SELECT bp.block_hash FROM blocktx.block_processing bp
		LEFT JOIN blocktx.blocks b ON b.hash = bp.block_hash
		WHERE b.processed_at IS NULL AND bp.processed_by = $1;
	`

	rows, err := p.db.QueryContext(ctx, q, processedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hashes := make([]*chainhash.Hash, 0)

	for rows.Next() {
		var hash []byte

		err = rows.Scan(&hash)
		if err != nil {
			return nil, err
		}

		txHash, err := chainhash.NewHash(hash)
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, txHash)
	}

	return hashes, nil
}
