package postgresql

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/bitcoin-sv/arc/internal/testdata"

	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"github.com/bitcoin-sv/arc/pkg/blocktx/blocktx_api"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Block struct {
	ID            int64     `db:"id"`
	Hash          string    `db:"hash"`
	PreviousHash  string    `db:"prevhash"`
	MerkleRoot    string    `db:"merkleroot"`
	MerklePath    *string   `db:"merkle_path"`
	Height        int64     `db:"height"`
	Orphaned      bool      `db:"orphanedyn"`
	Size          *int64    `db:"size"`
	TxCount       *int64    `db:"tx_count"`
	Processed     bool      `db:"processed"`
	ProcessedAt   time.Time `db:"processed_at"`
	InsertedAt    time.Time `db:"inserted_at"`
	InsertedAtNum int64     `db:"inserted_at_num"`
}

type Transaction struct {
	ID         int64  `db:"id"`
	Hash       []byte `db:"hash"`
	MerklePath string `db:"merkle_path"`
}

type BlockTransactionMap struct {
	BlockID       int64 `db:"blockid"`
	TransactionID int64 `db:"txid"`
	Pos           int64 `db:"pos"`
}

const (
	postgresPort   = "5432"
	migrationsPath = "file://migrations"
	dbName         = "main_test"
	dbUsername     = "arcuser"
	dbPassword     = "arcpass"
)

var dbInfo string

func revChainhash(t *testing.T, hashString string) *chainhash.Hash {
	hash, err := hex.DecodeString(hashString)
	require.NoError(t, err)
	txHash, err := chainhash.NewHash(hash)
	require.NoError(t, err)

	return txHash
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	port := "5434"
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15.4",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
			fmt.Sprintf("POSTGRES_USER=%s", dbUsername),
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			postgresPort: {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
		config.Tmpfs = map[string]string{
			"/var/lib/postgresql/data": "",
		}
	})
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	hostPort := resource.GetPort("5432/tcp")

	dbInfo = fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", hostPort, dbUsername, dbPassword, dbName)
	fmt.Println(dbInfo)
	var postgresDB *PostgreSQL
	err = pool.Retry(func() error {
		postgresDB, err = New(dbInfo, 10, 10)
		if err != nil {
			return err
		}
		return postgresDB.db.Ping()
	})
	if err != nil {
		log.Fatalf("failed to connect to docker: %s", err)
	}

	driver, err := migratepostgres.WithInstance(postgresDB.db, &migratepostgres.Config{
		MigrationsTable: "blocktx",
	})
	if err != nil {
		log.Fatalf("failed to create driver: %v", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("failed to initialize migrate instance: %v", err)
	}
	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to initialize migrate instance: %v", err)
	}
	code := m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf("failed to purge pool: %v", err)
	}

	os.Exit(code)
}

func loadFixtures(db *sql.DB, path string) error {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgresql"),
		testfixtures.Directory(path), // The directory containing the YAML files
	)
	if err != nil {
		log.Fatalf("failed to create fixtures: %v", err)
	}

	err = fixtures.Load()
	if err != nil {
		return fmt.Errorf("failed to load fixtures: %v", err)
	}

	return nil
}

func pruneTables(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE blocks;")
	if err != nil {
		return err
	}

	_, err = db.Exec("TRUNCATE TABLE transactions;")
	if err != nil {
		return err
	}

	_, err = db.Exec("TRUNCATE TABLE block_transactions_map;")
	if err != nil {
		return err
	}

	return nil
}

func TestPostgresDB(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	now := time.Date(2023, 12, 22, 12, 0, 0, 0, time.UTC)
	postgresDB, err := New(dbInfo, 10, 10, WithNow(func() time.Time {
		return now
	}))
	ctx := context.Background()
	require.NoError(t, err)
	defer func() {
		postgresDB.Close()
	}()

	t.Run("insert block / get block", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		blockHash1 := revChainhash(t, "000000000000000001b8adefc1eb98896c80e30e517b9e2655f1f929d9958a48")
		blockHash2 := revChainhash(t, "00000000000000000a081a539601645abe977946f8f6466a3c9e0c34d50be4a8")
		merkleRoot := revChainhash(t, "31e25c5ac7c143687f55fc49caf0f552ba6a16d4f785e4c9a9a842179a085f0c")
		block := &blocktx_api.Block{
			Hash:         blockHash2[:],
			PreviousHash: blockHash1[:],
			MerkleRoot:   merkleRoot[:],
			Height:       100,
		}

		id, err := postgresDB.InsertBlock(ctx, block)
		require.NoError(t, err)
		require.Equal(t, uint64(1), id)

		blockResp, err := postgresDB.GetBlock(ctx, blockHash2)
		require.NoError(t, err)
		require.Equal(t, block, blockResp)
	})

	t.Run("register transactions", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		txs := []*blocktx_api.TransactionAndSource{
			{
				Hash: testdata.TX1Hash[:],
			},
			{
				Hash: testdata.TX2Hash[:],
			},
			{
				Hash: testdata.TX3Hash[:],
			},
			{
				Hash: testdata.TX4Hash[:],
			},
		}
		err = postgresDB.RegisterTransactions(context.Background(), txs)
		require.NoError(t, err)

		d, err := sqlx.Open("postgres", dbInfo)
		require.NoError(t, err)

		var storedtx Transaction
		err = d.Get(&storedtx, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", string(txs[0].GetHash()))
		require.NoError(t, err)

		require.True(t, bytes.Equal(txs[0].GetHash(), storedtx.Hash))

		err = d.Get(&storedtx, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", string(txs[1].GetHash()))
		require.NoError(t, err)

		require.True(t, bytes.Equal(txs[1].GetHash(), storedtx.Hash))

		err = d.Get(&storedtx, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", string(txs[2].GetHash()))
		require.NoError(t, err)

		require.True(t, bytes.Equal(txs[2].GetHash(), storedtx.Hash))

		err = d.Get(&storedtx, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", string(txs[3].GetHash()))
		require.NoError(t, err)

		require.True(t, bytes.Equal(txs[3].GetHash(), storedtx.Hash))

		newHash1, err := hex.DecodeString("bb6817dad4c9a0e34803c2f1dec13ac8b89d539a57749be808d61f314bbbd855")
		require.NoError(t, err)

		newHash2, err := hex.DecodeString("ec52c6f777e6cb254d596001050f4151ef74cea78b6d50163407bf54bdd54e3f")
		require.NoError(t, err)
		// Register transactions which are partly already registered
		txs2 := []*blocktx_api.TransactionAndSource{
			txs[0],
			txs[2],
			{
				Hash: newHash1,
			},
			{
				Hash: newHash2,
			},
		}
		err = postgresDB.RegisterTransactions(context.Background(), txs2)
		require.NoError(t, err)

		err = d.Get(&storedtx, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", string(txs2[2].GetHash()))
		require.NoError(t, err)

		require.True(t, bytes.Equal(txs2[2].GetHash(), storedtx.Hash))

		err = d.Get(&storedtx, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", string(txs2[3].GetHash()))
		require.NoError(t, err)

		require.True(t, bytes.Equal(txs2[3].GetHash(), storedtx.Hash))
	})

	t.Run("get block gaps", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		require.NoError(t, loadFixtures(postgresDB.db, "fixtures/get_block_gaps"))

		blockGaps, err := postgresDB.GetBlockGaps(ctx, 12)
		require.NoError(t, err)
		require.Equal(t, 4, len(blockGaps))

		hash822014 := revChainhash(t, "67708796ef57464ed9eaf2a663d3da32372e4c2fb65558020000000000000000")
		hash822019 := revChainhash(t, "5696fc6e504b6aa2ae5d9c46b9418192dc61bd1b2e3364030000000000000000")
		hash822020 := revChainhash(t, "76404890880cb36ce68100abb05b3a958e17c0ed274d5c0a0000000000000000")
		hash822009 := revChainhash(t, "4ad773b1a464129a0ed8c7a8c71bb98175f0f01da1793f0e0000000000000000")

		expectedBlockGaps := []*store.BlockGap{
			{ // gap
				Height: 822019,
				Hash:   hash822019,
			},
			{ // gap
				Height: 822014,
				Hash:   hash822014,
			},
			{ // processing not finished
				Height: 822020,
				Hash:   hash822020,
			},
			{ // gap
				Height: 822009,
				Hash:   hash822009,
			},
		}

		require.ElementsMatch(t, expectedBlockGaps, blockGaps)
	})

	t.Run("update block transactions", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))
		postgresDB.maxPostgresBulkInsertRows = 5

		require.NoError(t, loadFixtures(postgresDB.db, "fixtures/update_block_transactions"))

		testMerklePaths := []string{"test1", "test2", "test3"}
		testBlockID := uint64(9736)

		txHash1 := revChainhash(t, "76732b80598326a18d3bf0a86518adbdf95d0ddc6ff6693004440f4776168c3b")
		txHash2 := revChainhash(t, "164e85a5d5bc2b2372e8feaa266e5e4b7d0808f8d2b784fb1f7349c4726392b0")

		txHashNotRegistered := revChainhash(t, "edd33fdcdfa68444d227780e2b62a4437c00120c5320d2026aeb24a781f4c3f1")

		_, err = postgresDB.UpdateBlockTransactions(context.Background(), testBlockID, []*blocktx_api.TransactionAndSource{
			{
				Hash: txHash1[:],
			},
			{
				Hash: txHash2[:],
			},
			{
				Hash: txHashNotRegistered[:],
			},
		}, []string{"test1"})

		require.ErrorContains(t, err, "transactions (len=3) and Merkle paths (len=1) have not the same lengths")

		updateResult, err := postgresDB.UpdateBlockTransactions(context.Background(), testBlockID, []*blocktx_api.TransactionAndSource{
			{
				Hash: txHash1[:],
			},
			{
				Hash: txHash2[:],
			},
			{
				Hash: txHashNotRegistered[:],
			},
		}, testMerklePaths)
		require.NoError(t, err)

		require.True(t, bytes.Equal(txHash1[:], updateResult[0].TxHash))
		require.Equal(t, testMerklePaths[0], updateResult[0].MerklePath)

		require.True(t, bytes.Equal(txHash2[:], updateResult[1].TxHash))
		require.Equal(t, testMerklePaths[1], updateResult[1].MerklePath)

		d, err := sqlx.Open("postgres", dbInfo)
		require.NoError(t, err)

		var storedtx Transaction

		err = d.Get(&storedtx, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", txHash1[:])
		require.NoError(t, err)

		require.Equal(t, txHash1[:], storedtx.Hash)
		require.Equal(t, testMerklePaths[0], storedtx.MerklePath)

		var mp BlockTransactionMap
		err = d.Get(&mp, "SELECT blockid, txid, pos from block_transactions_map WHERE txid=$1", storedtx.ID)
		require.NoError(t, err)

		require.Equal(t, storedtx.ID, mp.TransactionID)
		require.Equal(t, testBlockID, uint64(mp.BlockID))

		var storedtx2 Transaction

		err = d.Get(&storedtx2, "SELECT id, hash, merkle_path from transactions WHERE hash=$1", txHash2[:])
		require.NoError(t, err)

		require.Equal(t, txHash2[:], storedtx2.Hash)
		require.Equal(t, testMerklePaths[1], storedtx2.MerklePath)

		var mp2 BlockTransactionMap
		err = d.Get(&mp2, "SELECT blockid, txid, pos from block_transactions_map WHERE txid=$1", storedtx2.ID)
		require.NoError(t, err)

		require.Equal(t, storedtx2.ID, mp2.TransactionID)
		require.Equal(t, testBlockID, uint64(mp2.BlockID))

		// update exceeds max batch size

		txHash3 := revChainhash(t, "b4201cc6fc5768abff14adf75042ace6061da9176ee5bb943291b9ba7d7f5743")
		txHash4 := revChainhash(t, "37bd6c87927e75faeb3b3c939f64721cda48e1bb98742676eebe83aceee1a669")
		txHash5 := revChainhash(t, "952f80e20a0330f3b9c2dfd1586960064e797218b5c5df665cada221452c17eb")
		txHash6 := revChainhash(t, "861a281b27de016e50887288de87eab5ca56a1bb172cdff6dba965474ce0f608")
		txHash7 := revChainhash(t, "9421cc760c5405af950a76dc3e4345eaefd4e7322f172a3aee5e0ddc7b4f8313")
		txHash8 := revChainhash(t, "8b7d038db4518ac4c665abfc5aeaacbd2124ad8ca70daa8465ed2c4427c41b9b")

		_, err = postgresDB.UpdateBlockTransactions(context.Background(), testBlockID, []*blocktx_api.TransactionAndSource{
			{
				Hash: txHash3[:],
			},
			{
				Hash: txHash4[:],
			},
			{
				Hash: txHash5[:],
			},
			{
				Hash: txHash6[:],
			},
			{
				Hash: txHash7[:],
			},
			{
				Hash: txHash8[:],
			},
		}, []string{"test1", "test2", "test3", "test4", "test5", "test6"})
		require.NoError(t, err)
	})

	t.Run("test getting mined txs", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		require.NoError(t, loadFixtures(postgresDB.db, "fixtures/get_mined_transactions"))

		txHash1 := revChainhash(t, "76732b80598326a18d3bf0a86518adbdf95d0ddc6ff6693004440f4776168c3b")
		txHash2 := revChainhash(t, "164e85a5d5bc2b2372e8feaa266e5e4b7d0808f8d2b784fb1f7349c4726392b0")
		txHash3 := revChainhash(t, "dbbd24251b9bb824566412395bb76a579bca3477c2d0b4cbc210a769d3bb4177")
		txHash4 := revChainhash(t, "0d60dd6dc1f2649efb2847f801dfaa61361a438deb526da2de5b6875e0016514")

		blockHash := revChainhash(t, "6258b02da70a3e367e4c993b049fa9b76ef8f090ef9fd2010000000000000000")

		// get mined transaction and corresponding block
		minedTxs, err := postgresDB.GetMinedTransactions(ctx, []*chainhash.Hash{txHash1, txHash2, txHash3, txHash4})
		require.NoError(t, err)

		require.Len(t, minedTxs, 3)

		for _, tx := range minedTxs {
			require.True(t, bytes.Equal(tx.TxHash, txHash2[:]) || bytes.Equal(tx.TxHash, txHash3[:]) || bytes.Equal(tx.TxHash, txHash4[:]))
			require.Equal(t, tx.BlockHash, blockHash[:])
			require.Equal(t, uint64(826481), tx.BlockHeight)
		}
	})

	t.Run("clear data", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		require.NoError(t, loadFixtures(postgresDB.db, "fixtures/clear_data"))

		resp, err := postgresDB.ClearBlocktxTable(context.Background(), 10, "blocks")
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.Rows)

		d, err := sqlx.Open("postgres", dbInfo)
		require.NoError(t, err)
		var blocks []Block

		require.NoError(t, d.Select(&blocks, "SELECT id FROM blocks"))

		require.Len(t, blocks, 1)

		resp, err = postgresDB.ClearBlocktxTable(context.Background(), 10, "block_transactions_map")
		require.NoError(t, err)
		require.Equal(t, int64(5), resp.Rows)

		var mps []BlockTransactionMap
		require.NoError(t, d.Select(&mps, "SELECT blockid FROM block_transactions_map"))

		require.Len(t, mps, 5)

		resp, err = postgresDB.ClearBlocktxTable(context.Background(), 10, "transactions")
		require.NoError(t, err)
		require.Equal(t, int64(5), resp.Rows)

		var txs []Transaction
		require.NoError(t, d.Select(&txs, "SELECT id FROM transactions"))

		require.Len(t, txs, 5)
	})

	t.Run("set/get/del block processing", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		require.NoError(t, loadFixtures(postgresDB.db, "fixtures/block_processing"))

		bh1 := revChainhash(t, "747468cf7e6639ba9aa277ade1cf27639b0f214cec5719020000000000000000")

		processedBy, err := postgresDB.SetBlockProcessing(ctx, bh1, "pod-1")
		require.NoError(t, err)
		require.Equal(t, "pod-1", processedBy)

		// set a second time, expect error
		processedBy, err = postgresDB.SetBlockProcessing(ctx, bh1, "pod-1")
		require.ErrorIs(t, err, store.ErrBlockProcessingDuplicateKey)
		require.Equal(t, "pod-1", processedBy)

		bhInProgress := revChainhash(t, "f97e20396f02ab990ed31b9aec70c240f48b7e5ea239aa050000000000000000")

		blockHashes, err := postgresDB.GetBlockHashesProcessingInProgress(ctx, "pod-2")
		require.NoError(t, err)
		require.Len(t, blockHashes, 1)
		require.True(t, bhInProgress.IsEqual(blockHashes[0]))

		err = postgresDB.DelBlockProcessing(ctx, bhInProgress, "pod-1")
		require.ErrorIs(t, err, store.ErrBlockNotFound)

		err = postgresDB.DelBlockProcessing(ctx, bhInProgress, "pod-2")
		require.NoError(t, err)

		blockHashes, err = postgresDB.GetBlockHashesProcessingInProgress(ctx, "pod-2")
		require.NoError(t, err)
		require.Len(t, blockHashes, 0)
	})

	t.Run("mark block as done", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		require.NoError(t, loadFixtures(postgresDB.db, "fixtures/mark_block_as_done"))

		bh1 := revChainhash(t, "b71ab063c5f96cad71cdc59dcc94182a20a69cbd7eed2d070000000000000000")

		err = postgresDB.MarkBlockAsDone(ctx, bh1, 500, 75)
		require.NoError(t, err)

		d, err := sqlx.Open("postgres", dbInfo)
		require.NoError(t, err)
		var block Block

		require.NoError(t, d.Get(&block, "SELECT * FROM blocks WHERE hash=$1", bh1[:]))

		require.NotNil(t, block.Size)
		require.Equal(t, int64(500), *block.Size)

		require.NotNil(t, block.TxCount)
		require.Equal(t, int64(75), *block.TxCount)
		require.Equal(t, now.UTC(), block.ProcessedAt.UTC())
	})

	t.Run("verify merkle roots", func(t *testing.T) {
		defer require.NoError(t, pruneTables(postgresDB.db))

		require.NoError(t, loadFixtures(postgresDB.db, "fixtures/verify_merkle_roots"))

		merkleRequests := []*blocktx_api.MerkleRootVerificationRequest{
			{
				// correct merkleroot - should not return its height
				MerkleRoot:  "80d8a6a8306626fa25bd7e1bdae587805fc1f7d41e8aba0353488f8094153d4f",
				BlockHeight: 812010,
			},
			{
				// correct merkleroot but incorrect height - should return its height
				MerkleRoot:  "80d8a6a8306626fa25bd7e1bdae587805fc1f7d41e8aba0353488f8094153d4f",
				BlockHeight: 812011,
			},
			{
				// incorrect merkleroot - should return its height
				MerkleRoot:  "incorrect_merkle_root",
				BlockHeight: 822010,
			},
			{
				// merkleroot below lowest height - should not return its height
				MerkleRoot:  "fc377968f64fdcb6386ca7fcb6e6f8693a988663a07f552269823099909ea790",
				BlockHeight: 800000,
			},
			{
				// merkleroot above top height, but within limits - should not return its height
				MerkleRoot:  "fc377968f64fdcb6386ca7fcb6e6f8693a988663a07f552269823099909ea790",
				BlockHeight: 822030,
			},
			{
				// merkleroot above top height and above limits - should return its height
				MerkleRoot:  "fc377968f64fdcb6386ca7fcb6e6f8693a988663a07f552269823099909ea790",
				BlockHeight: 822032,
			},
		}
		maxAllowedBlockHeightMismatch := 10
		expectedUnverifiedBlockHeights := []uint64{812011, 822010, 822032}

		res, err := postgresDB.VerifyMerkleRoots(ctx, merkleRequests, maxAllowedBlockHeightMismatch)
		require.NoError(t, err)

		assert.Equal(t, expectedUnverifiedBlockHeights, res.UnverifiedBlockHeights)
	})
}
