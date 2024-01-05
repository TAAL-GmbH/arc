// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	context "context"
	"github.com/bitcoin-sv/arc/blocktx/blocktx_api"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	sync "sync"
)

// Ensure, that BlockTxAPIClientMock does implement blocktx_api.BlockTxAPIClient.
// If this is not the case, regenerate this file with moq.
var _ blocktx_api.BlockTxAPIClient = &BlockTxAPIClientMock{}

// BlockTxAPIClientMock is a mock implementation of blocktx_api.BlockTxAPIClient.
//
//	func TestSomethingThatUsesBlockTxAPIClient(t *testing.T) {
//
//		// make and configure a mocked blocktx_api.BlockTxAPIClient
//		mockedBlockTxAPIClient := &BlockTxAPIClientMock{
//			GetBlockFunc: func(ctx context.Context, in *blocktx_api.Hash, opts ...grpc.CallOption) (*blocktx_api.Block, error) {
//				panic("mock out the GetBlock method")
//			},
//			GetBlockForHeightFunc: func(ctx context.Context, in *blocktx_api.Height, opts ...grpc.CallOption) (*blocktx_api.Block, error) {
//				panic("mock out the GetBlockForHeight method")
//			},
//			GetBlockNotificationStreamFunc: func(ctx context.Context, in *blocktx_api.Height, opts ...grpc.CallOption) (blocktx_api.BlockTxAPI_GetBlockNotificationStreamClient, error) {
//				panic("mock out the GetBlockNotificationStream method")
//			},
//			GetBlockTransactionsFunc: func(ctx context.Context, in *blocktx_api.Block, opts ...grpc.CallOption) (*blocktx_api.Transactions, error) {
//				panic("mock out the GetBlockTransactions method")
//			},
//			GetLastProcessedBlockFunc: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*blocktx_api.Block, error) {
//				panic("mock out the GetLastProcessedBlock method")
//			},
//			GetMinedTransactionsForBlockFunc: func(ctx context.Context, in *blocktx_api.BlockAndSource, opts ...grpc.CallOption) (*blocktx_api.MinedTransactions, error) {
//				panic("mock out the GetMinedTransactionsForBlock method")
//			},
//			GetTransactionBlocksFunc: func(ctx context.Context, in *blocktx_api.Transactions, opts ...grpc.CallOption) (*blocktx_api.TransactionBlocks, error) {
//				panic("mock out the GetTransactionBlocks method")
//			},
//			GetTransactionMerklePathFunc: func(ctx context.Context, in *blocktx_api.Transaction, opts ...grpc.CallOption) (*blocktx_api.MerklePath, error) {
//				panic("mock out the GetTransactionMerklePath method")
//			},
//			HealthFunc: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*blocktx_api.HealthResponse, error) {
//				panic("mock out the Health method")
//			},
//			RegisterTransactionFunc: func(ctx context.Context, in *blocktx_api.TransactionAndSource, opts ...grpc.CallOption) (*blocktx_api.RegisterTransactionResponse, error) {
//				panic("mock out the RegisterTransaction method")
//			},
//		}
//
//		// use mockedBlockTxAPIClient in code that requires blocktx_api.BlockTxAPIClient
//		// and then make assertions.
//
//	}
type BlockTxAPIClientMock struct {
	// GetBlockFunc mocks the GetBlock method.
	GetBlockFunc func(ctx context.Context, in *blocktx_api.Hash, opts ...grpc.CallOption) (*blocktx_api.Block, error)

	// GetBlockForHeightFunc mocks the GetBlockForHeight method.
	GetBlockForHeightFunc func(ctx context.Context, in *blocktx_api.Height, opts ...grpc.CallOption) (*blocktx_api.Block, error)

	// GetBlockNotificationStreamFunc mocks the GetBlockNotificationStream method.
	GetBlockNotificationStreamFunc func(ctx context.Context, in *blocktx_api.Height, opts ...grpc.CallOption) (blocktx_api.BlockTxAPI_GetBlockNotificationStreamClient, error)

	// GetBlockTransactionsFunc mocks the GetBlockTransactions method.
	GetBlockTransactionsFunc func(ctx context.Context, in *blocktx_api.Block, opts ...grpc.CallOption) (*blocktx_api.Transactions, error)

	// GetLastProcessedBlockFunc mocks the GetLastProcessedBlock method.
	GetLastProcessedBlockFunc func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*blocktx_api.Block, error)

	// GetMinedTransactionsForBlockFunc mocks the GetMinedTransactionsForBlock method.
	GetMinedTransactionsForBlockFunc func(ctx context.Context, in *blocktx_api.BlockAndSource, opts ...grpc.CallOption) (*blocktx_api.MinedTransactions, error)

	// GetTransactionBlocksFunc mocks the GetTransactionBlocks method.
	GetTransactionBlocksFunc func(ctx context.Context, in *blocktx_api.Transactions, opts ...grpc.CallOption) (*blocktx_api.TransactionBlocks, error)

	// GetTransactionMerklePathFunc mocks the GetTransactionMerklePath method.
	GetTransactionMerklePathFunc func(ctx context.Context, in *blocktx_api.Transaction, opts ...grpc.CallOption) (*blocktx_api.MerklePath, error)

	// HealthFunc mocks the Health method.
	HealthFunc func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*blocktx_api.HealthResponse, error)

	// RegisterTransactionFunc mocks the RegisterTransaction method.
	RegisterTransactionFunc func(ctx context.Context, in *blocktx_api.TransactionAndSource, opts ...grpc.CallOption) (*blocktx_api.RegisterTransactionResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetBlock holds details about calls to the GetBlock method.
		GetBlock []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.Hash
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// GetBlockForHeight holds details about calls to the GetBlockForHeight method.
		GetBlockForHeight []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.Height
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// GetBlockNotificationStream holds details about calls to the GetBlockNotificationStream method.
		GetBlockNotificationStream []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.Height
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// GetBlockTransactions holds details about calls to the GetBlockTransactions method.
		GetBlockTransactions []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.Block
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// GetLastProcessedBlock holds details about calls to the GetLastProcessedBlock method.
		GetLastProcessedBlock []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *emptypb.Empty
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// GetMinedTransactionsForBlock holds details about calls to the GetMinedTransactionsForBlock method.
		GetMinedTransactionsForBlock []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.BlockAndSource
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// GetTransactionBlocks holds details about calls to the GetTransactionBlocks method.
		GetTransactionBlocks []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.Transactions
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// GetTransactionMerklePath holds details about calls to the GetTransactionMerklePath method.
		GetTransactionMerklePath []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.Transaction
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// Health holds details about calls to the Health method.
		Health []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *emptypb.Empty
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
		// RegisterTransaction holds details about calls to the RegisterTransaction method.
		RegisterTransaction []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *blocktx_api.TransactionAndSource
			// Opts is the opts argument value.
			Opts []grpc.CallOption
		}
	}
	lockGetBlock                     sync.RWMutex
	lockGetBlockForHeight            sync.RWMutex
	lockGetBlockNotificationStream   sync.RWMutex
	lockGetBlockTransactions         sync.RWMutex
	lockGetLastProcessedBlock        sync.RWMutex
	lockGetMinedTransactionsForBlock sync.RWMutex
	lockGetTransactionBlocks         sync.RWMutex
	lockGetTransactionMerklePath     sync.RWMutex
	lockHealth                       sync.RWMutex
	lockRegisterTransaction          sync.RWMutex
}

// GetBlock calls GetBlockFunc.
func (mock *BlockTxAPIClientMock) GetBlock(ctx context.Context, in *blocktx_api.Hash, opts ...grpc.CallOption) (*blocktx_api.Block, error) {
	if mock.GetBlockFunc == nil {
		panic("BlockTxAPIClientMock.GetBlockFunc: method is nil but BlockTxAPIClient.GetBlock was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.Hash
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetBlock.Lock()
	mock.calls.GetBlock = append(mock.calls.GetBlock, callInfo)
	mock.lockGetBlock.Unlock()
	return mock.GetBlockFunc(ctx, in, opts...)
}

// GetBlockCalls gets all the calls that were made to GetBlock.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetBlockCalls())
func (mock *BlockTxAPIClientMock) GetBlockCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.Hash
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.Hash
		Opts []grpc.CallOption
	}
	mock.lockGetBlock.RLock()
	calls = mock.calls.GetBlock
	mock.lockGetBlock.RUnlock()
	return calls
}

// GetBlockForHeight calls GetBlockForHeightFunc.
func (mock *BlockTxAPIClientMock) GetBlockForHeight(ctx context.Context, in *blocktx_api.Height, opts ...grpc.CallOption) (*blocktx_api.Block, error) {
	if mock.GetBlockForHeightFunc == nil {
		panic("BlockTxAPIClientMock.GetBlockForHeightFunc: method is nil but BlockTxAPIClient.GetBlockForHeight was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.Height
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetBlockForHeight.Lock()
	mock.calls.GetBlockForHeight = append(mock.calls.GetBlockForHeight, callInfo)
	mock.lockGetBlockForHeight.Unlock()
	return mock.GetBlockForHeightFunc(ctx, in, opts...)
}

// GetBlockForHeightCalls gets all the calls that were made to GetBlockForHeight.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetBlockForHeightCalls())
func (mock *BlockTxAPIClientMock) GetBlockForHeightCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.Height
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.Height
		Opts []grpc.CallOption
	}
	mock.lockGetBlockForHeight.RLock()
	calls = mock.calls.GetBlockForHeight
	mock.lockGetBlockForHeight.RUnlock()
	return calls
}

// GetBlockNotificationStream calls GetBlockNotificationStreamFunc.
func (mock *BlockTxAPIClientMock) GetBlockNotificationStream(ctx context.Context, in *blocktx_api.Height, opts ...grpc.CallOption) (blocktx_api.BlockTxAPI_GetBlockNotificationStreamClient, error) {
	if mock.GetBlockNotificationStreamFunc == nil {
		panic("BlockTxAPIClientMock.GetBlockNotificationStreamFunc: method is nil but BlockTxAPIClient.GetBlockNotificationStream was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.Height
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetBlockNotificationStream.Lock()
	mock.calls.GetBlockNotificationStream = append(mock.calls.GetBlockNotificationStream, callInfo)
	mock.lockGetBlockNotificationStream.Unlock()
	return mock.GetBlockNotificationStreamFunc(ctx, in, opts...)
}

// GetBlockNotificationStreamCalls gets all the calls that were made to GetBlockNotificationStream.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetBlockNotificationStreamCalls())
func (mock *BlockTxAPIClientMock) GetBlockNotificationStreamCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.Height
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.Height
		Opts []grpc.CallOption
	}
	mock.lockGetBlockNotificationStream.RLock()
	calls = mock.calls.GetBlockNotificationStream
	mock.lockGetBlockNotificationStream.RUnlock()
	return calls
}

// GetBlockTransactions calls GetBlockTransactionsFunc.
func (mock *BlockTxAPIClientMock) GetBlockTransactions(ctx context.Context, in *blocktx_api.Block, opts ...grpc.CallOption) (*blocktx_api.Transactions, error) {
	if mock.GetBlockTransactionsFunc == nil {
		panic("BlockTxAPIClientMock.GetBlockTransactionsFunc: method is nil but BlockTxAPIClient.GetBlockTransactions was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.Block
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetBlockTransactions.Lock()
	mock.calls.GetBlockTransactions = append(mock.calls.GetBlockTransactions, callInfo)
	mock.lockGetBlockTransactions.Unlock()
	return mock.GetBlockTransactionsFunc(ctx, in, opts...)
}

// GetBlockTransactionsCalls gets all the calls that were made to GetBlockTransactions.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetBlockTransactionsCalls())
func (mock *BlockTxAPIClientMock) GetBlockTransactionsCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.Block
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.Block
		Opts []grpc.CallOption
	}
	mock.lockGetBlockTransactions.RLock()
	calls = mock.calls.GetBlockTransactions
	mock.lockGetBlockTransactions.RUnlock()
	return calls
}

// GetLastProcessedBlock calls GetLastProcessedBlockFunc.
func (mock *BlockTxAPIClientMock) GetLastProcessedBlock(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*blocktx_api.Block, error) {
	if mock.GetLastProcessedBlockFunc == nil {
		panic("BlockTxAPIClientMock.GetLastProcessedBlockFunc: method is nil but BlockTxAPIClient.GetLastProcessedBlock was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *emptypb.Empty
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetLastProcessedBlock.Lock()
	mock.calls.GetLastProcessedBlock = append(mock.calls.GetLastProcessedBlock, callInfo)
	mock.lockGetLastProcessedBlock.Unlock()
	return mock.GetLastProcessedBlockFunc(ctx, in, opts...)
}

// GetLastProcessedBlockCalls gets all the calls that were made to GetLastProcessedBlock.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetLastProcessedBlockCalls())
func (mock *BlockTxAPIClientMock) GetLastProcessedBlockCalls() []struct {
	Ctx  context.Context
	In   *emptypb.Empty
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *emptypb.Empty
		Opts []grpc.CallOption
	}
	mock.lockGetLastProcessedBlock.RLock()
	calls = mock.calls.GetLastProcessedBlock
	mock.lockGetLastProcessedBlock.RUnlock()
	return calls
}

// GetMinedTransactionsForBlock calls GetMinedTransactionsForBlockFunc.
func (mock *BlockTxAPIClientMock) GetMinedTransactionsForBlock(ctx context.Context, in *blocktx_api.BlockAndSource, opts ...grpc.CallOption) (*blocktx_api.MinedTransactions, error) {
	if mock.GetMinedTransactionsForBlockFunc == nil {
		panic("BlockTxAPIClientMock.GetMinedTransactionsForBlockFunc: method is nil but BlockTxAPIClient.GetMinedTransactionsForBlock was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.BlockAndSource
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetMinedTransactionsForBlock.Lock()
	mock.calls.GetMinedTransactionsForBlock = append(mock.calls.GetMinedTransactionsForBlock, callInfo)
	mock.lockGetMinedTransactionsForBlock.Unlock()
	return mock.GetMinedTransactionsForBlockFunc(ctx, in, opts...)
}

// GetMinedTransactionsForBlockCalls gets all the calls that were made to GetMinedTransactionsForBlock.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetMinedTransactionsForBlockCalls())
func (mock *BlockTxAPIClientMock) GetMinedTransactionsForBlockCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.BlockAndSource
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.BlockAndSource
		Opts []grpc.CallOption
	}
	mock.lockGetMinedTransactionsForBlock.RLock()
	calls = mock.calls.GetMinedTransactionsForBlock
	mock.lockGetMinedTransactionsForBlock.RUnlock()
	return calls
}

// GetTransactionBlocks calls GetTransactionBlocksFunc.
func (mock *BlockTxAPIClientMock) GetTransactionBlocks(ctx context.Context, in *blocktx_api.Transactions, opts ...grpc.CallOption) (*blocktx_api.TransactionBlocks, error) {
	if mock.GetTransactionBlocksFunc == nil {
		panic("BlockTxAPIClientMock.GetTransactionBlocksFunc: method is nil but BlockTxAPIClient.GetTransactionBlocks was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.Transactions
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetTransactionBlocks.Lock()
	mock.calls.GetTransactionBlocks = append(mock.calls.GetTransactionBlocks, callInfo)
	mock.lockGetTransactionBlocks.Unlock()
	return mock.GetTransactionBlocksFunc(ctx, in, opts...)
}

// GetTransactionBlocksCalls gets all the calls that were made to GetTransactionBlocks.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetTransactionBlocksCalls())
func (mock *BlockTxAPIClientMock) GetTransactionBlocksCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.Transactions
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.Transactions
		Opts []grpc.CallOption
	}
	mock.lockGetTransactionBlocks.RLock()
	calls = mock.calls.GetTransactionBlocks
	mock.lockGetTransactionBlocks.RUnlock()
	return calls
}

// GetTransactionMerklePath calls GetTransactionMerklePathFunc.
func (mock *BlockTxAPIClientMock) GetTransactionMerklePath(ctx context.Context, in *blocktx_api.Transaction, opts ...grpc.CallOption) (*blocktx_api.MerklePath, error) {
	if mock.GetTransactionMerklePathFunc == nil {
		panic("BlockTxAPIClientMock.GetTransactionMerklePathFunc: method is nil but BlockTxAPIClient.GetTransactionMerklePath was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.Transaction
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockGetTransactionMerklePath.Lock()
	mock.calls.GetTransactionMerklePath = append(mock.calls.GetTransactionMerklePath, callInfo)
	mock.lockGetTransactionMerklePath.Unlock()
	return mock.GetTransactionMerklePathFunc(ctx, in, opts...)
}

// GetTransactionMerklePathCalls gets all the calls that were made to GetTransactionMerklePath.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.GetTransactionMerklePathCalls())
func (mock *BlockTxAPIClientMock) GetTransactionMerklePathCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.Transaction
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.Transaction
		Opts []grpc.CallOption
	}
	mock.lockGetTransactionMerklePath.RLock()
	calls = mock.calls.GetTransactionMerklePath
	mock.lockGetTransactionMerklePath.RUnlock()
	return calls
}

// Health calls HealthFunc.
func (mock *BlockTxAPIClientMock) Health(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*blocktx_api.HealthResponse, error) {
	if mock.HealthFunc == nil {
		panic("BlockTxAPIClientMock.HealthFunc: method is nil but BlockTxAPIClient.Health was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *emptypb.Empty
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockHealth.Lock()
	mock.calls.Health = append(mock.calls.Health, callInfo)
	mock.lockHealth.Unlock()
	return mock.HealthFunc(ctx, in, opts...)
}

// HealthCalls gets all the calls that were made to Health.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.HealthCalls())
func (mock *BlockTxAPIClientMock) HealthCalls() []struct {
	Ctx  context.Context
	In   *emptypb.Empty
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *emptypb.Empty
		Opts []grpc.CallOption
	}
	mock.lockHealth.RLock()
	calls = mock.calls.Health
	mock.lockHealth.RUnlock()
	return calls
}

// RegisterTransaction calls RegisterTransactionFunc.
func (mock *BlockTxAPIClientMock) RegisterTransaction(ctx context.Context, in *blocktx_api.TransactionAndSource, opts ...grpc.CallOption) (*blocktx_api.RegisterTransactionResponse, error) {
	if mock.RegisterTransactionFunc == nil {
		panic("BlockTxAPIClientMock.RegisterTransactionFunc: method is nil but BlockTxAPIClient.RegisterTransaction was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		In   *blocktx_api.TransactionAndSource
		Opts []grpc.CallOption
	}{
		Ctx:  ctx,
		In:   in,
		Opts: opts,
	}
	mock.lockRegisterTransaction.Lock()
	mock.calls.RegisterTransaction = append(mock.calls.RegisterTransaction, callInfo)
	mock.lockRegisterTransaction.Unlock()
	return mock.RegisterTransactionFunc(ctx, in, opts...)
}

// RegisterTransactionCalls gets all the calls that were made to RegisterTransaction.
// Check the length with:
//
//	len(mockedBlockTxAPIClient.RegisterTransactionCalls())
func (mock *BlockTxAPIClientMock) RegisterTransactionCalls() []struct {
	Ctx  context.Context
	In   *blocktx_api.TransactionAndSource
	Opts []grpc.CallOption
} {
	var calls []struct {
		Ctx  context.Context
		In   *blocktx_api.TransactionAndSource
		Opts []grpc.CallOption
	}
	mock.lockRegisterTransaction.RLock()
	calls = mock.calls.RegisterTransaction
	mock.lockRegisterTransaction.RUnlock()
	return calls
}
