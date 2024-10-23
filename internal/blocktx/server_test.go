package blocktx_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/bitcoin-sv/arc/internal/blocktx"
	"github.com/bitcoin-sv/arc/internal/blocktx/blocktx_api"
	"github.com/bitcoin-sv/arc/internal/blocktx/mocks"
	storeMocks "github.com/bitcoin-sv/arc/internal/blocktx/store/mocks"
	"github.com/bitcoin-sv/arc/internal/testdata"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/stretchr/testify/require"
)

func TestListenAndServe(t *testing.T) {
	tt := []struct {
		name string
	}{
		{
			name: "start and shutdown",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// given
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			storeMock := &storeMocks.BlocktxStoreMock{}
			pm := &mocks.PeerManagerMock{ShutdownFunc: func() {}}

			sut, err := blocktx.NewServer("", 0, logger, storeMock, pm, 0)
			require.NoError(t, err)
			defer sut.GracefulStop()

			// when
			err = sut.ListenAndServe("localhost:7000")

			// then
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		})
	}
}

func TestDelUnfinishedBlock(t *testing.T) {
	tt := []struct {
		name                                  string
		getBlockHashesProcessingInProgressErr error
		delBlockProcessingErr                 error

		expectedRows     int64
		expectedErrorStr string
	}{
		{
			name: "success",

			expectedRows: 6,
		},
		{
			name:                                  "error - getBlockHashesProcessingInProgress",
			getBlockHashesProcessingInProgressErr: errors.New("failed to get block hashes processing in progress"),

			expectedErrorStr: "failed to get block hashes processing in progress",
			expectedRows:     0,
		},
		{
			name:                  "error - delBlockProcessingErr",
			delBlockProcessingErr: errors.New("failed to delete block processing error"),

			expectedErrorStr: "failed to delete block processing error",
			expectedRows:     0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// given
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			storeMock := &storeMocks.BlocktxStoreMock{
				GetBlockHashesProcessingInProgressFunc: func(ctx context.Context, processedBy string) ([]*chainhash.Hash, error) {
					return []*chainhash.Hash{testdata.TX1Hash, testdata.TX2Hash}, tc.getBlockHashesProcessingInProgressErr
				},

				DelBlockProcessingFunc: func(ctx context.Context, hash *chainhash.Hash, processedBy string) (int64, error) {
					return 3, tc.delBlockProcessingErr
				},
			}

			sut, err := blocktx.NewServer("", 0, logger, storeMock, nil, 0)
			require.NoError(t, err)
			defer sut.GracefulStop()

			// when
			resp, err := sut.DelUnfinishedBlockProcessing(context.Background(), &blocktx_api.DelUnfinishedBlockProcessingRequest{
				ProcessedBy: "host",
			})

			// then
			if tc.expectedErrorStr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErrorStr)
				return
			}

			require.Equal(t, tc.expectedRows, resp.Rows)
		})
	}
}
