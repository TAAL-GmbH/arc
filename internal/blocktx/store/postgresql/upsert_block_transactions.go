package postgresql

import (
	"context"
	"errors"

	"github.com/lib/pq"

	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"github.com/bitcoin-sv/arc/internal/tracing"
)

// UpsertBlockTransactions upserts the transaction hashes for a given block hash.
func (p *PostgreSQL) UpsertBlockTransactions(ctx context.Context, blockID uint64, txsWithMerklePaths []store.TxWithMerklePath) error {
	ctx, span := tracing.StartTracing(ctx, "UpdateBlockTransactions", p.tracingEnabled, p.tracingAttributes...)
	defer tracing.EndTracing(span)

	txHashesBytes := make([][]byte, len(txsWithMerklePaths))
	merklePaths := make([]string, len(txsWithMerklePaths))
	for i, tx := range txsWithMerklePaths {
		txHashesBytes[i] = tx.Hash
		merklePaths[i] = tx.MerklePath
	}

	qUpsertTransactions := `
		WITH inserted_transactions AS (
				INSERT INTO blocktx.transactions (hash)
				SELECT UNNEST($2::BYTEA[])
				ON CONFLICT (hash)
				DO UPDATE SET hash = EXCLUDED.hash
				RETURNING id, hash
		)

		INSERT INTO blocktx.block_transactions_map (blockid, txid, merkle_path)
		SELECT
				$1::BIGINT,
				it.id,
				t.merkle_path
		FROM inserted_transactions it
		JOIN LATERAL UNNEST($2::BYTEA[], $3::TEXT[]) AS t(hash, merkle_path) ON it.hash = t.hash;
	`

	_, err := p.db.ExecContext(ctx, qUpsertTransactions, blockID, pq.Array(txHashesBytes), pq.Array(merklePaths))
	if err != nil {
		return errors.Join(store.ErrFailedToExecuteTxUpdateQuery, err)
	}

	return nil
}
