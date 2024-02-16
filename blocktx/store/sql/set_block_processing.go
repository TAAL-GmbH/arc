package sql

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/arc/blocktx/store"
	"github.com/lib/pq"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
)

func (s *SQL) SetBlockProcessing(ctx context.Context, hash *chainhash.Hash, setProcessedBy string) (string, error) {

	// Try to set a block as being processed by this instance
	qInsert := `
		INSERT INTO block_processing (block_hash, processed_by)
		VALUES ($1 ,$2 )
		RETURNING processed_by
	`

	var processedBy string
	err := s.db.QueryRowContext(ctx, qInsert, hash[:], setProcessedBy).Scan(&processedBy)
	if err != nil {
		var pqErr *pq.Error
		errors.As(err, &pqErr)
		if pqErr.Code == pq.ErrorCode("23505") {
			return processedBy, store.ErrBlockProcessingDuplicateKey
		}

		return "", fmt.Errorf("failed to set block processing: %v", err)
	}

	return processedBy, nil
}

func (s *SQL) DelBlockProcessing(ctx context.Context, hash *chainhash.Hash, processedBy string) error {
	q := `
		DELETE FROM block_processing WHERE hash = $1 AND processed_by = $2;
    `

	_, err := s.db.ExecContext(ctx, q, hash[:], processedBy)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQL) GetBlockHashesProcessingInProgress(ctx context.Context, processedBy string) ([]*chainhash.Hash, error) {
	// Check how many blocks this instance is currently processing
	q := `
	SELECT bp.block_hash FROM block_processing bp
	LEFT JOIN blocks b ON b.hash = bp.block_hash
	WHERE b.processed_at IS NULL AND bp.processed_by = $1;
    `

	rows, err := s.db.QueryContext(ctx, q, processedBy)
	if err != nil {
		return nil, err
	}

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