package sql

import (
	"database/sql"

	"github.com/TAAL-GmbH/arc/blocktx/blocktx_api"

	"context"
)

// InsertBlockTransactions inserts the transaction hashes for a given block hash
func (s *SQL) InsertBlockTransactions(ctx context.Context, blockId uint64, transactions []*blocktx_api.Transaction) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	qTx := `
			INSERT INTO transactions (hash) VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING id
			;
		`
	qMap := `
		INSERT INTO block_transactions_map (
		 blockid
		,txid
		) VALUES (
		 $1
		,$2
		)
		ON CONFLICT DO NOTHING
	`

	// txn, err := s.db.BeginTx(ctx, nil)
	// if err != nil {
	// 	return err
	// }

	for _, tx := range transactions {
		var txid uint64

		if err := s.db.QueryRowContext(ctx, qTx, tx.Hash).Scan(&txid); err != nil {
			if err == sql.ErrNoRows {
				if err := s.db.QueryRowContext(ctx, "SELECT id FROM transactions WHERE hash = $1", tx.Hash).Scan(&txid); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		_, err := s.db.ExecContext(ctx, qMap, blockId, txid)
		if err != nil {
			return err
		}
	}

	// if err := txn.Commit(); err != nil {
	// 	return err
	// }

	return nil
}
