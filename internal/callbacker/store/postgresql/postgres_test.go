package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/bitcoin-sv/arc/internal/callbacker/store"
	"github.com/bitcoin-sv/arc/internal/testdata"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	tutils "github.com/bitcoin-sv/arc/internal/callbacker/store/postgresql/internal/tests"
)

const (
	postgresPort   = "5432"
	migrationsPath = "file://migrations"
	dbName         = "main_test"
	dbUsername     = "arcuser"
	dbPassword     = "arcpass"
)

var dbInfo string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	port := "5433"
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
	dbConn, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("failed to create db connection: %v", err)
	}
	err = pool.Retry(func() error {
		return dbConn.Ping()
	})

	if err != nil {
		log.Fatalf("failed to connect to docker: %s", err)
	}

	driver, err := migratepostgres.WithInstance(dbConn, &migratepostgres.Config{
		MigrationsTable: "callbacker",
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

func TestPostgresDBt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	now := time.Date(2023, 10, 1, 14, 25, 0, 0, time.UTC)

	postgresDB, err := New(dbInfo, 10, 10)
	require.NoError(t, err)
	defer postgresDB.Close()

	t.Run("set many", func(t *testing.T) {
		// given
		defer pruneTables(t, postgresDB.db)

		data := []*store.CallbackData{
			{
				Url:       "https://test-callback-1/",
				Token:     "token",
				TxID:      testdata.TX2,
				TxStatus:  "SEEN_ON_NETWORK",
				Timestamp: now,
			},
			{
				Url:         "https://test-callback-1/",
				Token:       "token",
				TxID:        testdata.TX2,
				TxStatus:    "MINED",
				Timestamp:   now,
				BlockHash:   ptrTo(testdata.Block1),
				BlockHeight: ptrTo(uint64(4524235)),
			},
			{
				// duplicate
				Url:         "https://test-callback-1/",
				Token:       "token",
				TxID:        testdata.TX2,
				TxStatus:    "MINED",
				Timestamp:   now,
				BlockHash:   &testdata.Block1,
				BlockHeight: ptrTo(uint64(4524235)),
			},
			{
				Url:          "https://test-callback-2/",
				Token:        "token",
				TxID:         testdata.TX3,
				TxStatus:     "MINED",
				Timestamp:    now,
				BlockHash:    &testdata.Block2,
				BlockHeight:  ptrTo(uint64(4524236)),
				CompetingTxs: []string{testdata.TX2},
			},
			{
				// duplicate
				Url:          "https://test-callback-2/",
				Token:        "token",
				TxID:         testdata.TX3,
				TxStatus:     "MINED",
				Timestamp:    now,
				BlockHash:    &testdata.Block2,
				BlockHeight:  ptrTo(uint64(4524236)),
				CompetingTxs: []string{testdata.TX2},
			},
			{
				Url:         "https://test-callback-2/",
				TxID:        testdata.TX2,
				TxStatus:    "MINED",
				Timestamp:   now,
				BlockHash:   &testdata.Block1,
				BlockHeight: ptrTo(uint64(4524235)),
			},
		}

		// when
		err = postgresDB.SetMany(context.Background(), data)

		// then
		require.NoError(t, err)

		// check if every UNIQUE record has been saved
		uniqueRecords := make(map[*store.CallbackData]struct{})

		for _, dr := range data {
			duplicate := false

			for ur := range uniqueRecords {
				if tutils.CallbackRecordEqual(ur, dr) {
					duplicate = true
					break
				}
			}

			if !duplicate {
				uniqueRecords[dr] = struct{}{}
			}
		}
		require.NotEmpty(t, uniqueRecords)
		fmt.Println(uniqueRecords)

		// read all from db
		dbCallbacks := tutils.ReadAllCallbacks(t, postgresDB.db)
		for _, c := range dbCallbacks {
			// check if exists in unique records
			for ur := range uniqueRecords {
				if ur == nil {
					continue
				}

				if tutils.CallbackRecordEqual(ur, c) {
					// remove
					delete(uniqueRecords, ur)
					break
				}
			}
		}

		require.Empty(t, uniqueRecords)
	})

	t.Run("pop many", func(t *testing.T) {
		// given
		defer pruneTables(t, postgresDB.db)
		loadFixtures(t, postgresDB.db, "fixtures/pop_many")

		const concurentCalls = 5
		const popLimit = 10

		// count current records
		count := tutils.CountCallbacks(t, postgresDB.db)
		require.GreaterOrEqual(t, count, concurentCalls*popLimit)

		ctx := context.Background()
		start := make(chan struct{})
		rm := sync.Map{}
		wg := sync.WaitGroup{}

		// when
		wg.Add(concurentCalls)
		for i := range concurentCalls {
			go func() {
				defer wg.Done()
				<-start

				records, err := postgresDB.PopMany(ctx, popLimit)
				require.NoError(t, err)

				rm.Store(i, records)
			}()
		}

		close(start) // signal all goroutines to start
		wg.Wait()

		// then
		count2 := tutils.CountCallbacks(t, postgresDB.db)
		require.Equal(t, count-concurentCalls*popLimit, count2)

		for i := range concurentCalls {
			records, ok := rm.Load(i)
			require.True(t, ok)

			callbacks := records.([]*store.CallbackData)
			require.Equal(t, popLimit, len(callbacks))
		}
	})

}

func pruneTables(t *testing.T, db *sql.DB) error {
	t.Helper()

	_, err := db.Exec("TRUNCATE TABLE callbacker.callbacks;")
	if err != nil {
		return err
	}

	return nil
}

func loadFixtures(t *testing.T, db *sql.DB, path string) {
	t.Helper()

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgresql"),
		testfixtures.Directory(path), // The directory containing the YAML files
	)
	if err != nil {
		t.Fatalf("failed to create fixtures: %v", err)
	}

	err = fixtures.Load()
	if err != nil {
		t.Fatalf("failed to load fixtures: %v", err)
	}
}
