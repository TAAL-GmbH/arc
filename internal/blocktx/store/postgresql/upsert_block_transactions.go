package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"

	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"github.com/bitcoin-sv/arc/internal/tracing"
)

// UpsertBlockTransactions upserts the transaction hashes for a given block hash and returns updated registered transactions hashes.
func (p *PostgreSQL) UpsertBlockTransactions(ctx context.Context, blockID uint64, txsWithMerklePaths []store.TxWithMerklePath) (err error) {
	ctx, span := tracing.StartTracing(ctx, "UpsertBlockTransactions", p.tracingEnabled, append(p.tracingAttributes, attribute.Int("updates", len(txsWithMerklePaths)))...)
	defer func() {
		tracing.EndTracing(span, err)
	}()

	txHashes := make([][]byte, len(txsWithMerklePaths))
	blockIDs := make([]uint64, len(txsWithMerklePaths))
	merkleTreeIndexes := make([]int64, len(txsWithMerklePaths))
	for pos, tx := range txsWithMerklePaths {
		txHashes[pos] = tx.Hash
		blockIDs[pos] = blockID
		merkleTreeIndexes[pos] = tx.MerkleTreeIndex
	}

	var copyRows [][]any

	qBulkUpsert := `
		INSERT INTO blocktx.transactions (hash)
			SELECT UNNEST($1::BYTEA[])
			ON CONFLICT (hash)
			DO UPDATE SET hash = EXCLUDED.hash
		RETURNING id, is_registered`

	_, spanTxs := tracing.StartTracing(ctx, "UpsertBlockTransactions_transactions", p.tracingEnabled, append(p.tracingAttributes, attribute.Int("updates", len(txsWithMerklePaths)))...)
	rows, err := p.db.QueryContext(ctx, qBulkUpsert, pq.Array(txHashes))
	if err != nil {
		return errors.Join(store.ErrFailedToUpsertTransactions, err)
	}
	tracing.EndTracing(spanTxs, nil)

	counter := 0
	txIDs := make([]uint64, len(txsWithMerklePaths))
	for rows.Next() {
		var txID uint64
		var isRegistered bool

		err = rows.Scan(&txID, &isRegistered)
		if err != nil {
			return errors.Join(store.ErrFailedToGetRows, err)
		}

		txIDs[counter] = txID

		mp := ""
		if isRegistered {
			mp = txsWithMerklePaths[counter].MerklePath
		}

		copyRows = append(copyRows, []any{blockID, txID, mp, merkleTreeIndexes[counter]})

		counter++
	}

	if len(txIDs) != len(txsWithMerklePaths) {
		return errors.Join(store.ErrMismatchedTxIDsAndMerklePathLength, err)
	}

	conn, err := pgx.Connect(ctx, p.dbInfo)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, spanTxsMap := tracing.StartTracing(ctx, "UpsertBlockTransactions_transactions_map", p.tracingEnabled, append(p.tracingAttributes, attribute.Int("updates", len(txsWithMerklePaths)))...)

	_, err = conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"blocktx", "block_transactions_map"},
		[]string{"blockid", "txid", "merkle_path", "merkle_tree_index"},
		pgx.CopyFromRows(copyRows),
	)
	if err != nil {
		return errors.Join(store.ErrFailedToUpsertBlockTransactionsMap, err)
	}
	tracing.EndTracing(spanTxsMap, nil)

	return nil
}
