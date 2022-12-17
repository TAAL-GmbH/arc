package sql

import (
	"context"

	"github.com/TAAL-GmbH/arc/blocktx/blocktx_api"
)

// GetBlockTransactions returns the transaction hashes for a given block hash
func (s *SQL) GetBlockTransactions(ctx context.Context, block *blocktx_api.Block) (*blocktx_api.Transactions, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	q := `
		SELECT 
		 t.hash
		FROM transactions t
		INNER JOIN block_transactions_map m ON m.txid = t.id
		INNER JOIN blocks b ON m.blockid = b.id
		WHERE b.hash = $1
	`

	rows, err := s.db.QueryContext(ctx, q, block.Hash)
	if err != nil {
		return nil, err

	}

	defer rows.Close()

	var hash []byte
	var transactions []*blocktx_api.Transaction

	for rows.Next() {
		err = rows.Scan(&hash)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &blocktx_api.Transaction{
			Hash: hash,
		})
	}

	return &blocktx_api.Transactions{
		Transactions: transactions,
	}, nil
}
