package sql

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/arc/blocktx/blocktx_api"
	"github.com/bitcoin-sv/arc/blocktx/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransactionBlockSuite struct {
	DatabaseTestSuite
}

func (s GetTransactionBlockSuite) Test() {
	block := GetTestBlock()
	tx := GetTestTransaction()
	s.InsertBlock(block)

	s.InsertTransaction(tx)

	s.InsertBlockTransactionMap(&store.BlockTransactionMap{
		BlockID:       block.ID,
		TransactionID: tx.ID,
		Pos:           2,
	})

	store, err := NewPostgresStore(defaultParams)
	require.NoError(s.T(), err)

	b, err := store.GetTransactionBlock(context.Background(), &blocktx_api.Transaction{
		Hash: []byte(tx.Hash),
	})
	require.NoError(s.T(), err)

	assert.Equal(s.T(), block.Hash, string(b.Hash))
}

func TestGetTransactionBlock(t *testing.T) {
	s := new(GetTransactionBlockSuite)
	suite.Run(t, s)
	if err := recover(); err != nil {
		require.NoError(t, s.Database.Stop())
	}
}
