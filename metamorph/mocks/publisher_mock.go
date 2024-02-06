// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/bitcoin-sv/arc/metamorph/async"
	"sync"
)

// Ensure, that PublisherMock does implement async.Publisher.
// If this is not the case, regenerate this file with moq.
var _ async.Publisher = &PublisherMock{}

// PublisherMock is a mock implementation of async.Publisher.
//
//	func TestSomethingThatUsesPublisher(t *testing.T) {
//
//		// make and configure a mocked async.Publisher
//		mockedPublisher := &PublisherMock{
//			PublishTransactionFunc: func(ctx context.Context, hash []byte) error {
//				panic("mock out the PublishTransaction method")
//			},
//		}
//
//		// use mockedPublisher in code that requires async.Publisher
//		// and then make assertions.
//
//	}
type PublisherMock struct {
	// PublishTransactionFunc mocks the PublishTransaction method.
	PublishTransactionFunc func(ctx context.Context, hash []byte) error

	// calls tracks calls to the methods.
	calls struct {
		// PublishTransaction holds details about calls to the PublishTransaction method.
		PublishTransaction []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Hash is the hash argument value.
			Hash []byte
		}
	}
	lockPublishTransaction sync.RWMutex
}

// PublishTransaction calls PublishTransactionFunc.
func (mock *PublisherMock) PublishTransaction(ctx context.Context, hash []byte) error {
	if mock.PublishTransactionFunc == nil {
		panic("PublisherMock.PublishTransactionFunc: method is nil but Publisher.PublishTransaction was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Hash []byte
	}{
		Ctx:  ctx,
		Hash: hash,
	}
	mock.lockPublishTransaction.Lock()
	mock.calls.PublishTransaction = append(mock.calls.PublishTransaction, callInfo)
	mock.lockPublishTransaction.Unlock()
	return mock.PublishTransactionFunc(ctx, hash)
}

// PublishTransactionCalls gets all the calls that were made to PublishTransaction.
// Check the length with:
//
//	len(mockedPublisher.PublishTransactionCalls())
func (mock *PublisherMock) PublishTransactionCalls() []struct {
	Ctx  context.Context
	Hash []byte
} {
	var calls []struct {
		Ctx  context.Context
		Hash []byte
	}
	mock.lockPublishTransaction.RLock()
	calls = mock.calls.PublishTransaction
	mock.lockPublishTransaction.RUnlock()
	return calls
}
