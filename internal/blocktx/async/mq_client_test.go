package async

import (
	"context"
	"errors"
	"testing"

	"github.com/bitcoin-sv/arc/internal/testdata"
	"github.com/bitcoin-sv/arc/pkg/blocktx/blocktx_api"
	"github.com/stretchr/testify/require"
)

//go:generate moq -out ./nats_client_mock.go . NatsClient

func TestMQClient_PublishMinedTxs(t *testing.T) {
	txBlock := &blocktx_api.TransactionBlock{
		BlockHash:       testdata.Block1Hash[:],
		BlockHeight:     1,
		TransactionHash: testdata.TX1Hash[:],
		MerklePath:      "mp-1",
	}

	tt := []struct {
		name       string
		txsBlocks  []*blocktx_api.TransactionBlock
		publishErr error

		expectedErrorStr     string
		expectedPublishCalls int
	}{
		{
			name: "small batch",
			txsBlocks: []*blocktx_api.TransactionBlock{
				txBlock,
				txBlock,
			},

			expectedPublishCalls: 1,
		},
		{
			name: "exactly batch size",
			txsBlocks: []*blocktx_api.TransactionBlock{
				txBlock, txBlock, txBlock, txBlock, txBlock,
			},

			expectedPublishCalls: 1,
		},
		{
			name: "large batch",
			txsBlocks: []*blocktx_api.TransactionBlock{
				txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock,
				txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock,
				txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock, txBlock,
				txBlock,
			},

			expectedPublishCalls: 7,
		},
		{
			name: "publish err",
			txsBlocks: []*blocktx_api.TransactionBlock{
				txBlock, txBlock,
			},
			publishErr: errors.New("failed to publish"),

			expectedErrorStr:     "failed to publish",
			expectedPublishCalls: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			natsMock := &NatsClientMock{
				PublishFunc: func(subj string, data []byte) error {
					return tc.publishErr
				},
			}

			txChannel := make(chan []byte, 10)
			requestTxChannel := make(chan []byte, 10)

			mqClient := NewNatsMQClient(natsMock, txChannel, requestTxChannel, WithMaxBatchSize(5), WithTracer())

			err := mqClient.PublishMinedTxs(context.Background(), tc.txsBlocks)

			if tc.expectedErrorStr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErrorStr)
				return
			}

			require.Equal(t, tc.expectedPublishCalls, len(natsMock.PublishCalls()))

		})
	}
}