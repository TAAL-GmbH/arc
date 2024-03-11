// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/bitcoin-sv/arc/internal/metamorph"
	"sync"
)

// Ensure, that TransactionHandlerMock does implement metamorph.TransactionHandler.
// If this is not the case, regenerate this file with moq.
var _ metamorph.TransactionHandler = &TransactionHandlerMock{}

// TransactionHandlerMock is a mock implementation of metamorph.TransactionHandler.
//
//	func TestSomethingThatUsesTransactionHandler(t *testing.T) {
//
//		// make and configure a mocked metamorph.TransactionHandler
//		mockedTransactionHandler := &TransactionHandlerMock{
//			GetTransactionFunc: func(ctx context.Context, txID string) ([]byte, error) {
//				panic("mock out the GetTransaction method")
//			},
//			GetTransactionStatusFunc: func(ctx context.Context, txID string) (*metamorph.TransactionStatus, error) {
//				panic("mock out the GetTransactionStatus method")
//			},
//			HealthFunc: func(ctx context.Context) error {
//				panic("mock out the Health method")
//			},
//			SubmitTransactionFunc: func(ctx context.Context, tx []byte, options *metamorph.TransactionOptions) (*metamorph.TransactionStatus, error) {
//				panic("mock out the SubmitTransaction method")
//			},
//			SubmitTransactionsFunc: func(ctx context.Context, tx [][]byte, options *metamorph.TransactionOptions) ([]*metamorph.TransactionStatus, error) {
//				panic("mock out the SubmitTransactions method")
//			},
//		}
//
//		// use mockedTransactionHandler in code that requires metamorph.TransactionHandler
//		// and then make assertions.
//
//	}
type TransactionHandlerMock struct {
	// GetTransactionFunc mocks the GetTransaction method.
	GetTransactionFunc func(ctx context.Context, txID string) ([]byte, error)

	// GetTransactionStatusFunc mocks the GetTransactionStatus method.
	GetTransactionStatusFunc func(ctx context.Context, txID string) (*metamorph.TransactionStatus, error)

	// HealthFunc mocks the Health method.
	HealthFunc func(ctx context.Context) error

	// SubmitTransactionFunc mocks the SubmitTransaction method.
	SubmitTransactionFunc func(ctx context.Context, tx []byte, options *metamorph.TransactionOptions) (*metamorph.TransactionStatus, error)

	// SubmitTransactionsFunc mocks the SubmitTransactions method.
	SubmitTransactionsFunc func(ctx context.Context, tx [][]byte, options *metamorph.TransactionOptions) ([]*metamorph.TransactionStatus, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetTransaction holds details about calls to the GetTransaction method.
		GetTransaction []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// TxID is the txID argument value.
			TxID string
		}
		// GetTransactionStatus holds details about calls to the GetTransactionStatus method.
		GetTransactionStatus []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// TxID is the txID argument value.
			TxID string
		}
		// Health holds details about calls to the Health method.
		Health []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// SubmitTransaction holds details about calls to the SubmitTransaction method.
		SubmitTransaction []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Tx is the tx argument value.
			Tx []byte
			// Options is the options argument value.
			Options *metamorph.TransactionOptions
		}
		// SubmitTransactions holds details about calls to the SubmitTransactions method.
		SubmitTransactions []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Tx is the tx argument value.
			Tx [][]byte
			// Options is the options argument value.
			Options *metamorph.TransactionOptions
		}
	}
	lockGetTransaction       sync.RWMutex
	lockGetTransactionStatus sync.RWMutex
	lockHealth               sync.RWMutex
	lockSubmitTransaction    sync.RWMutex
	lockSubmitTransactions   sync.RWMutex
}

// GetTransaction calls GetTransactionFunc.
func (mock *TransactionHandlerMock) GetTransaction(ctx context.Context, txID string) ([]byte, error) {
	if mock.GetTransactionFunc == nil {
		panic("TransactionHandlerMock.GetTransactionFunc: method is nil but TransactionHandler.GetTransaction was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		TxID string
	}{
		Ctx:  ctx,
		TxID: txID,
	}
	mock.lockGetTransaction.Lock()
	mock.calls.GetTransaction = append(mock.calls.GetTransaction, callInfo)
	mock.lockGetTransaction.Unlock()
	return mock.GetTransactionFunc(ctx, txID)
}

// GetTransactionCalls gets all the calls that were made to GetTransaction.
// Check the length with:
//
//	len(mockedTransactionHandler.GetTransactionCalls())
func (mock *TransactionHandlerMock) GetTransactionCalls() []struct {
	Ctx  context.Context
	TxID string
} {
	var calls []struct {
		Ctx  context.Context
		TxID string
	}
	mock.lockGetTransaction.RLock()
	calls = mock.calls.GetTransaction
	mock.lockGetTransaction.RUnlock()
	return calls
}

// GetTransactionStatus calls GetTransactionStatusFunc.
func (mock *TransactionHandlerMock) GetTransactionStatus(ctx context.Context, txID string) (*metamorph.TransactionStatus, error) {
	if mock.GetTransactionStatusFunc == nil {
		panic("TransactionHandlerMock.GetTransactionStatusFunc: method is nil but TransactionHandler.GetTransactionStatus was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		TxID string
	}{
		Ctx:  ctx,
		TxID: txID,
	}
	mock.lockGetTransactionStatus.Lock()
	mock.calls.GetTransactionStatus = append(mock.calls.GetTransactionStatus, callInfo)
	mock.lockGetTransactionStatus.Unlock()
	return mock.GetTransactionStatusFunc(ctx, txID)
}

// GetTransactionStatusCalls gets all the calls that were made to GetTransactionStatus.
// Check the length with:
//
//	len(mockedTransactionHandler.GetTransactionStatusCalls())
func (mock *TransactionHandlerMock) GetTransactionStatusCalls() []struct {
	Ctx  context.Context
	TxID string
} {
	var calls []struct {
		Ctx  context.Context
		TxID string
	}
	mock.lockGetTransactionStatus.RLock()
	calls = mock.calls.GetTransactionStatus
	mock.lockGetTransactionStatus.RUnlock()
	return calls
}

// Health calls HealthFunc.
func (mock *TransactionHandlerMock) Health(ctx context.Context) error {
	if mock.HealthFunc == nil {
		panic("TransactionHandlerMock.HealthFunc: method is nil but TransactionHandler.Health was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockHealth.Lock()
	mock.calls.Health = append(mock.calls.Health, callInfo)
	mock.lockHealth.Unlock()
	return mock.HealthFunc(ctx)
}

// HealthCalls gets all the calls that were made to Health.
// Check the length with:
//
//	len(mockedTransactionHandler.HealthCalls())
func (mock *TransactionHandlerMock) HealthCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockHealth.RLock()
	calls = mock.calls.Health
	mock.lockHealth.RUnlock()
	return calls
}

// SubmitTransaction calls SubmitTransactionFunc.
func (mock *TransactionHandlerMock) SubmitTransaction(ctx context.Context, tx []byte, options *metamorph.TransactionOptions) (*metamorph.TransactionStatus, error) {
	if mock.SubmitTransactionFunc == nil {
		panic("TransactionHandlerMock.SubmitTransactionFunc: method is nil but TransactionHandler.SubmitTransaction was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Tx      []byte
		Options *metamorph.TransactionOptions
	}{
		Ctx:     ctx,
		Tx:      tx,
		Options: options,
	}
	mock.lockSubmitTransaction.Lock()
	mock.calls.SubmitTransaction = append(mock.calls.SubmitTransaction, callInfo)
	mock.lockSubmitTransaction.Unlock()
	return mock.SubmitTransactionFunc(ctx, tx, options)
}

// SubmitTransactionCalls gets all the calls that were made to SubmitTransaction.
// Check the length with:
//
//	len(mockedTransactionHandler.SubmitTransactionCalls())
func (mock *TransactionHandlerMock) SubmitTransactionCalls() []struct {
	Ctx     context.Context
	Tx      []byte
	Options *metamorph.TransactionOptions
} {
	var calls []struct {
		Ctx     context.Context
		Tx      []byte
		Options *metamorph.TransactionOptions
	}
	mock.lockSubmitTransaction.RLock()
	calls = mock.calls.SubmitTransaction
	mock.lockSubmitTransaction.RUnlock()
	return calls
}

// SubmitTransactions calls SubmitTransactionsFunc.
func (mock *TransactionHandlerMock) SubmitTransactions(ctx context.Context, tx [][]byte, options *metamorph.TransactionOptions) ([]*metamorph.TransactionStatus, error) {
	if mock.SubmitTransactionsFunc == nil {
		panic("TransactionHandlerMock.SubmitTransactionsFunc: method is nil but TransactionHandler.SubmitTransactions was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Tx      [][]byte
		Options *metamorph.TransactionOptions
	}{
		Ctx:     ctx,
		Tx:      tx,
		Options: options,
	}
	mock.lockSubmitTransactions.Lock()
	mock.calls.SubmitTransactions = append(mock.calls.SubmitTransactions, callInfo)
	mock.lockSubmitTransactions.Unlock()
	return mock.SubmitTransactionsFunc(ctx, tx, options)
}

// SubmitTransactionsCalls gets all the calls that were made to SubmitTransactions.
// Check the length with:
//
//	len(mockedTransactionHandler.SubmitTransactionsCalls())
func (mock *TransactionHandlerMock) SubmitTransactionsCalls() []struct {
	Ctx     context.Context
	Tx      [][]byte
	Options *metamorph.TransactionOptions
} {
	var calls []struct {
		Ctx     context.Context
		Tx      [][]byte
		Options *metamorph.TransactionOptions
	}
	mock.lockSubmitTransactions.RLock()
	calls = mock.calls.SubmitTransactions
	mock.lockSubmitTransactions.RUnlock()
	return calls
}
