package sql

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/arc/blocktx/blocktx_api"
	"github.com/bitcoin-sv/arc/blocktx/store"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InsertBlockSuite struct {
	DatabaseTestSuite
}

func (s InsertBlockSuite) Run() {
	block := GetTestBlock()

	pstore, err := NewPostgresStore(defaultParams)
	require.NoError(s.T(), err)

	b := blocktx_api.Block{
		Hash:         []byte(block.Hash),
		PreviousHash: []byte(block.PreviousHash),
		MerkleRoot:   []byte(block.MerkleRoot),
		Height:       uint64(block.Height),
		Orphaned:     block.Orphaned,
	}

	bid, err := pstore.InsertBlock(context.Background(), &b)
	require.NoError(s.T(), err)

	d, err := sqlx.Open("postgres", defaultParams.String())
	require.NoError(s.T(), err)

	var blk store.Block

	err = d.Get(&blk, "SELECT hash, prevhash, merkleroot, height, orphanedyn from blocks WHERE id=$1", bid)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), block.Hash, blk.Hash)
	assert.Equal(s.T(), block.PreviousHash, blk.PreviousHash)
	assert.Equal(s.T(), block.MerklePath, blk.MerklePath)
	assert.Equal(s.T(), block.MerkleRoot, blk.MerkleRoot)
}

func TestInsertBlockSuite(t *testing.T) {
	s := new(InsertBlockSuite)
	suite.Run(t, s)
	if err := recover(); err != nil {
		require.NoError(t, s.Database.Stop())
	}
}
