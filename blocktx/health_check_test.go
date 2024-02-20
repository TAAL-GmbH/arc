package blocktx_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/bitcoin-sv/arc/blocktx"
	"github.com/bitcoin-sv/arc/blocktx/mocks"
	"github.com/bitcoin-sv/arc/blocktx/store"
	"github.com/libsv/go-p2p/wire"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestCheck(t *testing.T) {
	tt := []struct {
		name               string
		service            string
		pingErr            error
		processorHealthErr error
		primaryErr         error

		expectedStatus grpc_health_v1.HealthCheckResponse_ServingStatus
	}{
		{
			name:       "liveness - peer not found",
			service:    "readiness",
			primaryErr: nil,

			expectedStatus: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		},
		{
			name:       "db error - not connected",
			service:    "readiness",
			primaryErr: errors.New("not connected"),

			expectedStatus: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		},
		{
			name:       "db error - not connected",
			service:    "readiness",
			primaryErr: nil,

			expectedStatus: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			const batchSize = 4

			var storeMock = &store.InterfaceMock{
				GetBlockGapsFunc: func(ctx context.Context, heightRange int) ([]*store.BlockGap, error) {
					return nil, nil
				},

				GetPrimaryFunc: func(ctx context.Context) (string, error) {
					return "", tc.primaryErr
				},
			}

			req := &grpc_health_v1.HealthCheckRequest{
				Service: tc.service,
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			peerHandler, err := blocktx.NewPeerHandler(logger, storeMock, 100, []string{}, wire.TestNet, blocktx.WithTransactionBatchSize(batchSize))
			require.NoError(t, err)

			server := blocktx.NewServer(storeMock, logger, peerHandler)
			resp, err := server.Check(context.Background(), req)
			require.NoError(t, err)

			require.Equal(t, tc.expectedStatus, resp.Status)

			peerHandler.Shutdown()
		})
	}
}

func TestWatch(t *testing.T) {
	tt := []struct {
		name       string
		service    string
		pingErr    error
		primaryErr error

		expectedStatus grpc_health_v1.HealthCheckResponse_ServingStatus
	}{
		{
			name:       "liveness - healthy",
			service:    "liveness",
			primaryErr: nil,

			expectedStatus: grpc_health_v1.HealthCheckResponse_SERVING,
		},
		{
			name:       "not ready - healthy",
			service:    "readiness",
			primaryErr: nil,

			expectedStatus: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		},
		{
			name:       "not ready - healthy",
			service:    "readiness",
			primaryErr: errors.New("not connected"),

			expectedStatus: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			const batchSize = 4

			var storeMock = &store.InterfaceMock{
				GetBlockGapsFunc: func(ctx context.Context, heightRange int) ([]*store.BlockGap, error) {
					return nil, nil
				},

				GetPrimaryFunc: func(ctx context.Context) (string, error) {
					return "", tc.primaryErr
				},
			}

			req := &grpc_health_v1.HealthCheckRequest{
				Service: tc.service,
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			peerHandler, err := blocktx.NewPeerHandler(logger, storeMock, 100, []string{}, wire.TestNet, blocktx.WithTransactionBatchSize(batchSize))
			require.NoError(t, err)

			server := blocktx.NewServer(storeMock, logger, peerHandler)

			watchServer := &mocks.HealthWatchServerMock{
				SendFunc: func(healthCheckResponse *grpc_health_v1.HealthCheckResponse) error {
					require.Equal(t, tc.expectedStatus, healthCheckResponse.Status)
					return nil
				},
			}

			err = server.Watch(req, watchServer)
			require.NoError(t, err)
			peerHandler.Shutdown()
		})
	}
}