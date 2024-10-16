package metamorph_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bitcoin-sv/arc/internal/metamorph/metamorph_api"

	apiMocks "github.com/bitcoin-sv/arc/internal/metamorph/mocks"
	"github.com/bitcoin-sv/arc/internal/testdata"
	"github.com/bitcoin-sv/arc/pkg/metamorph"
	"github.com/bitcoin-sv/arc/pkg/metamorph/mocks"

	sdkTx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestClient_SetUnlockedByName(t *testing.T) {
	tt := []struct {
		name           string
		setUnlockedErr error

		expectedErrorStr string
	}{
		{
			name: "success",
		},
		{
			name:           "err",
			setUnlockedErr: errors.New("failed to clear data"),

			expectedErrorStr: "failed to clear data",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			apiClient := &apiMocks.MetaMorphAPIClientMock{
				SetUnlockedByNameFunc: func(ctx context.Context, in *metamorph_api.SetUnlockedByNameRequest, opts ...grpc.CallOption) (*metamorph_api.SetUnlockedByNameResponse, error) {
					return &metamorph_api.SetUnlockedByNameResponse{RecordsAffected: 5}, tc.setUnlockedErr
				},
			}

			client := metamorph.NewClient(apiClient)
			// When
			res, err := client.SetUnlockedByName(context.Background(), "test-1")
			// Then
			if tc.expectedErrorStr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErrorStr)
				return
			}

			require.Equal(t, int64(5), res)
		})
	}
}

func TestClient_SubmitTransaction(t *testing.T) {
	now := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)
	tt := []struct {
		name               string
		options            *metamorph.TransactionOptions
		putTxErr           error
		putTxStatus        *metamorph_api.TransactionStatus
		withMqClient       bool
		publishSubmitTxErr error

		expectedErrorStr string
		expectedStatus   *metamorph.TransactionStatus
	}{
		{
			name: "wait for received",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_RECEIVED,
			},
			putTxStatus: &metamorph_api.TransactionStatus{
				Txid:   testdata.TX1Hash.String(),
				Status: metamorph_api.Status_RECEIVED,
			},

			expectedStatus: &metamorph.TransactionStatus{
				TxID:      testdata.TX1Hash.String(),
				Status:    metamorph_api.Status_RECEIVED.String(),
				Timestamp: now.Unix(),
			},
		},
		{
			name: "wait for seen - double spend attempted",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_SEEN_ON_NETWORK,
			},
			putTxStatus: &metamorph_api.TransactionStatus{
				Txid:         testdata.TX1Hash.String(),
				Status:       metamorph_api.Status_DOUBLE_SPEND_ATTEMPTED,
				CompetingTxs: []string{"1234"},
			},

			expectedStatus: &metamorph.TransactionStatus{
				TxID:         testdata.TX1Hash.String(),
				Status:       metamorph_api.Status_DOUBLE_SPEND_ATTEMPTED.String(),
				Timestamp:    now.Unix(),
				CompetingTxs: []string{"1234"},
			},
		},
		{
			name: "wait for received, put tx err, no mq client",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_RECEIVED,
			},
			putTxStatus: &metamorph_api.TransactionStatus{
				Txid:   testdata.TX1Hash.String(),
				Status: metamorph_api.Status_RECEIVED,
			},
			putTxErr:     errors.New("failed to put tx"),
			withMqClient: false,

			expectedStatus:   nil,
			expectedErrorStr: "failed to put tx",
		},
		{
			name: "wait for queued, with mq client",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_QUEUED,
			},
			withMqClient: true,

			expectedStatus: &metamorph.TransactionStatus{
				TxID:      testdata.TX1Hash.String(),
				Status:    metamorph_api.Status_QUEUED.String(),
				Timestamp: now.Unix(),
			},
		},
		{
			name: "wait for queued, with mq client, publish submit tx err",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_QUEUED,
			},
			withMqClient:       true,
			publishSubmitTxErr: errors.New("failed to publish tx"),

			expectedStatus:   nil,
			expectedErrorStr: "failed to publish tx",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			apiClient := &apiMocks.MetaMorphAPIClientMock{
				PutTransactionFunc: func(ctx context.Context, in *metamorph_api.TransactionRequest, opts ...grpc.CallOption) (*metamorph_api.TransactionStatus, error) {
					return tc.putTxStatus, tc.putTxErr
				},
			}

			opts := []func(client *metamorph.Metamorph){metamorph.WithNow(func() time.Time { return now })}
			if tc.withMqClient {
				mqClient := &mocks.MessageQueueClientMock{
					PublishMarshalFunc: func(topic string, m protoreflect.ProtoMessage) error { return tc.publishSubmitTxErr },
				}
				opts = append(opts, metamorph.WithMqClient(mqClient))
			}

			client := metamorph.NewClient(apiClient, opts...)
			// When
			tx, err := sdkTx.NewTransactionFromHex(testdata.TX1RawString)
			// Then
			require.NoError(t, err)
			status, err := client.SubmitTransaction(context.Background(), tx, tc.options)

			require.Equal(t, tc.expectedStatus, status)

			if tc.expectedErrorStr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErrorStr)
				return
			}
		})
	}
}

func TestClient_SubmitTransactions(t *testing.T) {
	now := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)
	tx1, err := sdkTx.NewTransactionFromHex("010000000000000000ef016c50da4e8941c9b11720a4a29b40955c30f246b25740cd1aecffa2e3c4acd144000000006b483045022100eaf7791ec8ec1b9766473e70a5e41ac1734b6e43126d3dfa142c5f7670256cae02206169a3d22f0519b2631e8b952d8530db3502f2494ba038672e69b23a1e03340c412102f87ce69f6ba5444aed49c34470041189c1e1060acd99341959c0594002c61bf0ffffffffbf070000000000001976a914c2b6fd4319122b9b5156a2a0060d19864c24f49a88ac01be070000000000001976a914c2b6fd4319122b9b5156a2a0060d19864c24f49a88ac00000000")
	require.NoError(t, err)
	tx2, err := sdkTx.NewTransactionFromHex("010000000000000000ef0159f09a1fc4f1df5790730de57f96840fc5fbbbb08ebff52c986fd43a842588e0000000006a473044022020152c7c9f09e6b31bce86fc2b21bf8b0e5edfdaba575196dc47b933b4ec6f9502201815515de957ff44d8f9a9368a055d04bc7e1a675ad8b34ca67e47b28459718f412102f87ce69f6ba5444aed49c34470041189c1e1060acd99341959c0594002c61bf0ffffffffbf070000000000001976a914c2b6fd4319122b9b5156a2a0060d19864c24f49a88ac01be070000000000001976a914c2b6fd4319122b9b5156a2a0060d19864c24f49a88ac00000000")
	require.NoError(t, err)
	tx3, err := sdkTx.NewTransactionFromHex("010000000000000000ef016a4c158eb2906c84b3d95206a4dac765baf4dff63120e09ab0134dc6505a23bf000000006b483045022100a46fb3431796212efc3f78b2a8559ba66a5f197977ee983b765a8a1497c0e31a022077ff0eed59beadbdd08b9806ed1f83e41e57ef4892564a75ecbeddd99f14f60f412102f87ce69f6ba5444aed49c34470041189c1e1060acd99341959c0594002c61bf0ffffffffbf070000000000001976a914c2b6fd4319122b9b5156a2a0060d19864c24f49a88ac01be070000000000001976a914c2b6fd4319122b9b5156a2a0060d19864c24f49a88ac00000000")
	require.NoError(t, err)
	tt := []struct {
		name               string
		options            *metamorph.TransactionOptions
		putTxErr           error
		putTxStatus        *metamorph_api.TransactionStatuses
		withMqClient       bool
		publishSubmitTxErr error

		expectedErrorStr string
		expectedStatuses []*metamorph.TransactionStatus
	}{
		{
			name: "wait for received",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_RECEIVED,
			},
			putTxStatus: &metamorph_api.TransactionStatuses{
				Statuses: []*metamorph_api.TransactionStatus{
					{
						Txid:   tx1.TxID(),
						Status: metamorph_api.Status_RECEIVED,
					},
					{
						Txid:   tx2.TxID(),
						Status: metamorph_api.Status_RECEIVED,
					},
					{
						Txid:   tx3.TxID(),
						Status: metamorph_api.Status_RECEIVED,
					},
				},
			},

			expectedStatuses: []*metamorph.TransactionStatus{
				{
					TxID:      tx1.TxID(),
					Status:    metamorph_api.Status_RECEIVED.String(),
					Timestamp: now.Unix(),
				},
				{
					TxID:      tx2.TxID(),
					Status:    metamorph_api.Status_RECEIVED.String(),
					Timestamp: now.Unix(),
				},
				{
					TxID:      tx3.TxID(),
					Status:    metamorph_api.Status_RECEIVED.String(),
					Timestamp: now.Unix(),
				},
			},
		},
		{
			name: "wait for received, put tx err, no mq client",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_RECEIVED,
			},
			putTxStatus: &metamorph_api.TransactionStatuses{
				Statuses: []*metamorph_api.TransactionStatus{
					{
						Txid:   tx1.TxID(),
						Status: metamorph_api.Status_RECEIVED,
					},
					{
						Txid:   tx2.TxID(),
						Status: metamorph_api.Status_RECEIVED,
					},
					{
						Txid:   tx3.TxID(),
						Status: metamorph_api.Status_RECEIVED,
					},
				},
			},
			putTxErr:     errors.New("failed to put tx"),
			withMqClient: false,

			expectedStatuses: nil,
			expectedErrorStr: "failed to put tx",
		},
		{
			name: "wait for queued, with mq client",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_QUEUED,
			},
			withMqClient: true,

			expectedStatuses: []*metamorph.TransactionStatus{
				{
					TxID:      tx1.TxID(),
					Status:    metamorph_api.Status_QUEUED.String(),
					Timestamp: now.Unix(),
				},
				{
					TxID:      tx2.TxID(),
					Status:    metamorph_api.Status_QUEUED.String(),
					Timestamp: now.Unix(),
				},
				{
					TxID:      tx3.TxID(),
					Status:    metamorph_api.Status_QUEUED.String(),
					Timestamp: now.Unix(),
				},
			},
		},
		{
			name: "wait for queued, with mq client, publish submit tx err",
			options: &metamorph.TransactionOptions{
				WaitForStatus: metamorph_api.Status_QUEUED,
			},
			withMqClient:       true,
			publishSubmitTxErr: errors.New("failed to publish tx"),

			expectedStatuses: nil,
			expectedErrorStr: "failed to publish tx",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			apiClient := &apiMocks.MetaMorphAPIClientMock{
				PutTransactionsFunc: func(ctx context.Context, in *metamorph_api.TransactionRequests, opts ...grpc.CallOption) (*metamorph_api.TransactionStatuses, error) {
					return tc.putTxStatus, tc.putTxErr
				},
			}

			opts := []func(client *metamorph.Metamorph){metamorph.WithNow(func() time.Time { return now })}
			if tc.withMqClient {
				mqClient := &mocks.MessageQueueClientMock{
					PublishMarshalFunc: func(topic string, m protoreflect.ProtoMessage) error {
						return tc.publishSubmitTxErr
					},
				}
				opts = append(opts, metamorph.WithMqClient(mqClient))
			}

			client := metamorph.NewClient(apiClient, opts...)
			// When
			statuses, err := client.SubmitTransactions(context.Background(), sdkTx.Transactions{tx1, tx2, tx3}, tc.options)
			// Then
			require.Equal(t, tc.expectedStatuses, statuses)

			if tc.expectedErrorStr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErrorStr)
				return
			}
		})
	}
}

func TestClient_GetTransaction(t *testing.T) {
	// Define test cases
	tt := []struct {
		name           string
		txID           string
		mockResp       *metamorph_api.Transaction
		mockErr        error
		expectedData   []byte
		expectedErrStr string
	}{
		{
			name:         "success - transaction found",
			txID:         "testTxID",
			mockResp:     &metamorph_api.Transaction{Txid: "testTxID", RawTx: []byte{0x01, 0x02, 0x03}},
			expectedData: []byte{0x01, 0x02, 0x03},
		},
		{
			name:           "error - transaction not found",
			txID:           "missingTxID",
			mockResp:       nil,
			mockErr:        metamorph.ErrTransactionNotFound,
			expectedErrStr: metamorph.ErrTransactionNotFound.Error(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockClient := &apiMocks.MetaMorphAPIClientMock{
				GetTransactionFunc: func(ctx context.Context, req *metamorph_api.TransactionStatusRequest, opts ...grpc.CallOption) (*metamorph_api.Transaction, error) {
					return tc.mockResp, tc.mockErr
				},
			}

			client := metamorph.NewClient(mockClient)

			// When
			txData, err := client.GetTransaction(context.Background(), tc.txID)

			// Then
			if tc.expectedErrStr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrStr)
				require.Nil(t, txData)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedData, txData)
			}

			// Ensure the mock was called once
			require.Equal(t, 1, len(mockClient.GetTransactionCalls()))
			require.Equal(t, tc.txID, mockClient.GetTransactionCalls()[0].In.Txid)
		})
	}
}

func TestClient_GetTransactions(t *testing.T) {
	tests := []struct {
		name             string
		txIDs            []string
		mockResponse     *metamorph_api.Transactions
		mockError        error
		expectedTxCount  int
		expectedErrorStr string
	}{
		{
			name:  "success - multiple transactions",
			txIDs: []string{"tx1", "tx2"},
			mockResponse: &metamorph_api.Transactions{
				Transactions: []*metamorph_api.Transaction{
					{
						Txid:        "tx1",
						RawTx:       []byte("transaction1"),
						BlockHeight: 100,
					},
					{
						Txid:        "tx2",
						RawTx:       []byte("transaction2"),
						BlockHeight: 101,
					},
				},
			},
			expectedTxCount: 2,
		},
		{
			name:             "error - transaction retrieval fails",
			txIDs:            []string{"tx1"},
			mockResponse:     nil,
			mockError:        errors.New("failed to retrieve transactions"),
			expectedErrorStr: "failed to retrieve transactions",
		},
		{
			name:             "error - transaction not found",
			txIDs:            []string{"tx1"},
			mockResponse:     nil,
			mockError:        errors.New("transaction not found"),
			expectedErrorStr: "transaction not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			apiClient := &apiMocks.MetaMorphAPIClientMock{
				GetTransactionsFunc: func(ctx context.Context, req *metamorph_api.TransactionsStatusRequest, opts ...grpc.CallOption) (*metamorph_api.Transactions, error) {
					require.Equal(t, tt.txIDs, req.TxIDs, "TxIDs do not match")
					return tt.mockResponse, tt.mockError
				},
			}

			client := metamorph.NewClient(apiClient)

			// When
			transactions, err := client.GetTransactions(context.Background(), tt.txIDs)

			// Then
			if tt.expectedErrorStr == "" {
				require.NoError(t, err)
				require.Len(t, transactions, tt.expectedTxCount)
				for i, tx := range transactions {
					require.Equal(t, tt.mockResponse.Transactions[i].Txid, tx.TxID, "Transaction ID mismatch")
					require.Equal(t, tt.mockResponse.Transactions[i].RawTx, tx.Bytes, "Raw transaction bytes mismatch")
					require.Equal(t, tt.mockResponse.Transactions[i].BlockHeight, tx.BlockHeight, "Block height mismatch")
				}
			} else {
				require.ErrorContains(t, err, tt.expectedErrorStr)
			}

		})
	}
}

func TestClient_GetTransactionStatus(t *testing.T) {
	tt := []struct {
		name           string
		txID           string
		status         metamorph_api.Status
		expectedStatus string
		expectedErrStr string
	}{
		{"unknown", "testTxID1", metamorph_api.Status_UNKNOWN, "UNKNOWN", "UNKNOWN"},
		{"queued", "testTxID2", metamorph_api.Status_QUEUED, "QUEUED", "QUEUED"},
		{"received", "testTxID3", metamorph_api.Status_RECEIVED, "RECEIVED", "RECEIVED"},
		{"stored", "testTxID4", metamorph_api.Status_STORED, "STORED", "STORED"},
		{"announced_to_network", "testTxID5", metamorph_api.Status_ANNOUNCED_TO_NETWORK, "ANNOUNCED_TO_NETWORK", "ANNOUNCED_TO_NETWORK"},
		{"requested_by_network", "testTxID6", metamorph_api.Status_REQUESTED_BY_NETWORK, "REQUESTED_BY_NETWORK", "REQUESTED_BY_NETWORK"},
		{"sent_to_network", "testTxID7", metamorph_api.Status_SENT_TO_NETWORK, "SENT_TO_NETWORK", "SENT_TO_NETWORK"},
		{"accepted_by_network", "testTxID8", metamorph_api.Status_ACCEPTED_BY_NETWORK, "ACCEPTED_BY_NETWORK", "ACCEPTED_BY_NETWORK"},
		{"seen_in_orphan_mempool", "testTxID9", metamorph_api.Status_SEEN_IN_ORPHAN_MEMPOOL, "SEEN_IN_ORPHAN_MEMPOOL", "SEEN_IN_ORPHAN_MEMPOOL"},
		{"seen_on_network", "testTxID10", metamorph_api.Status_SEEN_ON_NETWORK, "SEEN_ON_NETWORK", "SEEN_ON_NETWORK"},
		{"double_spend_attempted", "testTxID11", metamorph_api.Status_DOUBLE_SPEND_ATTEMPTED, "DOUBLE_SPEND_ATTEMPTED", "DOUBLE_SPEND_ATTEMPTED"},
		{"rejected", "testTxID12", metamorph_api.Status_REJECTED, "REJECTED", "REJECTED"},
		{"mined", "testTxID13", metamorph_api.Status_MINED, "MINED", "MINED"},
		{"transaction not found", "missingTxID", metamorph_api.Status_UNKNOWN, "", metamorph.ErrTransactionNotFound.Error()},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockClient := &apiMocks.MetaMorphAPIClientMock{
				GetTransactionStatusFunc: func(ctx context.Context, req *metamorph_api.TransactionStatusRequest, opts ...grpc.CallOption) (*metamorph_api.TransactionStatus, error) {
					if tc.expectedErrStr != "" {
						return nil, errors.New(tc.expectedErrStr)
					}
					return &metamorph_api.TransactionStatus{
						Txid:        tc.txID,
						Status:      tc.status,
						BlockHash:   "sampleBlockHash",
						BlockHeight: 100,
					}, nil
				},
			}

			now := time.Now().Unix()
			client := metamorph.NewClient(mockClient, metamorph.WithNow(func() time.Time { return time.Unix(now, 0) }))

			// When
			status, err := client.GetTransactionStatus(context.Background(), tc.txID)

			// Then
			if tc.expectedErrStr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrStr)
				require.Nil(t, status)
			} else {
				require.NoError(t, err)
				require.Equal(t, &metamorph.TransactionStatus{
					TxID:        tc.txID,
					Status:      tc.expectedStatus,
					BlockHash:   "sampleBlockHash",
					BlockHeight: 100,
					Timestamp:   now,
				}, status)
			}

			require.Equal(t, 1, len(mockClient.GetTransactionStatusCalls()), "Unexpected number of GetTransactionStatus calls")
		})
	}
}

func TestClient_Health(t *testing.T) {
	tt := []struct {
		name           string
		mockError      error
		expectedErrStr string
	}{
		{
			name:           "success - health check passes",
			mockError:      nil,
			expectedErrStr: "",
		},
		{
			name:           "error - health check fails",
			mockError:      errors.New("health check error"),
			expectedErrStr: "health check error",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockClient := &apiMocks.MetaMorphAPIClientMock{
				HealthFunc: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*metamorph_api.HealthResponse, error) {
					return &metamorph_api.HealthResponse{}, tc.mockError
				},
			}

			client := metamorph.NewClient(mockClient)

			// When
			err := client.Health(context.Background())

			// Then
			if tc.expectedErrStr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErrStr)
			}

			require.Equal(t, 1, len(mockClient.HealthCalls()), "Unexpected number of Health calls")
		})
	}
}

func TestGetTransactions(t *testing.T) {
	// Define test cases
	tt := []struct {
		name                 string
		txIDs                []string
		mockResp             *metamorph_api.Transactions
		mockErr              error
		expectedTransactions []*metamorph.Transaction
		expectedErrStr       string
	}{
		{
			name:  "success - valid transactions",
			txIDs: []string{"txid1", "txid2"},
			mockResp: &metamorph_api.Transactions{
				Transactions: []*metamorph_api.Transaction{
					{
						Txid:        "txid1",
						RawTx:       []byte{0x01, 0x02},
						BlockHeight: 100,
					},
					{
						Txid:        "txid2",
						RawTx:       []byte{0x03, 0x04},
						BlockHeight: 200,
					},
				},
			},
			expectedTransactions: []*metamorph.Transaction{
				{
					TxID:        "txid1",
					Bytes:       []byte{0x01, 0x02},
					BlockHeight: 100,
				},
				{
					TxID:        "txid2",
					Bytes:       []byte{0x03, 0x04},
					BlockHeight: 200,
				},
			},
		},
		{
			name:           "error - transaction not found",
			txIDs:          []string{"missingTxid"},
			mockResp:       nil,
			mockErr:        errors.New("transaction not found"),
			expectedErrStr: "transaction not found",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockClient := &apiMocks.MetaMorphAPIClientMock{
				GetTransactionsFunc: func(ctx context.Context, req *metamorph_api.TransactionsStatusRequest, opts ...grpc.CallOption) (*metamorph_api.Transactions, error) {
					return tc.mockResp, tc.mockErr
				},
			}

			client := metamorph.NewClient(mockClient)

			// When
			transactions, err := client.GetTransactions(context.Background(), tc.txIDs)

			// Then
			if tc.expectedErrStr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrStr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTransactions, transactions)
			}

		})
	}
}
