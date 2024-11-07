// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"sync"
)

// Ensure, that UnitOfWorkMock does implement store.UnitOfWork.
// If this is not the case, regenerate this file with moq.
var _ store.UnitOfWork = &UnitOfWorkMock{}

// UnitOfWorkMock is a mock implementation of store.UnitOfWork.
//
//	func TestSomethingThatUsesUnitOfWork(t *testing.T) {
//
//		// make and configure a mocked store.UnitOfWork
//		mockedUnitOfWork := &UnitOfWorkMock{
//			CommitFunc: func() error {
//				panic("mock out the Commit method")
//			},
//			RollbackFunc: func() error {
//				panic("mock out the Rollback method")
//			},
//			WriteLockBlocksTableFunc: func(ctx context.Context) error {
//				panic("mock out the WriteLockBlocksTable method")
//			},
//		}
//
//		// use mockedUnitOfWork in code that requires store.UnitOfWork
//		// and then make assertions.
//
//	}
type UnitOfWorkMock struct {
	// CommitFunc mocks the Commit method.
	CommitFunc func() error

	// RollbackFunc mocks the Rollback method.
	RollbackFunc func() error

	// WriteLockBlocksTableFunc mocks the WriteLockBlocksTable method.
	WriteLockBlocksTableFunc func(ctx context.Context) error

	// calls tracks calls to the methods.
	calls struct {
		// Commit holds details about calls to the Commit method.
		Commit []struct {
		}
		// Rollback holds details about calls to the Rollback method.
		Rollback []struct {
		}
		// WriteLockBlocksTable holds details about calls to the WriteLockBlocksTable method.
		WriteLockBlocksTable []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
	}
	lockCommit               sync.RWMutex
	lockRollback             sync.RWMutex
	lockWriteLockBlocksTable sync.RWMutex
}

// Commit calls CommitFunc.
func (mock *UnitOfWorkMock) Commit() error {
	if mock.CommitFunc == nil {
		panic("UnitOfWorkMock.CommitFunc: method is nil but UnitOfWork.Commit was just called")
	}
	callInfo := struct {
	}{}
	mock.lockCommit.Lock()
	mock.calls.Commit = append(mock.calls.Commit, callInfo)
	mock.lockCommit.Unlock()
	return mock.CommitFunc()
}

// CommitCalls gets all the calls that were made to Commit.
// Check the length with:
//
//	len(mockedUnitOfWork.CommitCalls())
func (mock *UnitOfWorkMock) CommitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockCommit.RLock()
	calls = mock.calls.Commit
	mock.lockCommit.RUnlock()
	return calls
}

// Rollback calls RollbackFunc.
func (mock *UnitOfWorkMock) Rollback() error {
	if mock.RollbackFunc == nil {
		panic("UnitOfWorkMock.RollbackFunc: method is nil but UnitOfWork.Rollback was just called")
	}
	callInfo := struct {
	}{}
	mock.lockRollback.Lock()
	mock.calls.Rollback = append(mock.calls.Rollback, callInfo)
	mock.lockRollback.Unlock()
	return mock.RollbackFunc()
}

// RollbackCalls gets all the calls that were made to Rollback.
// Check the length with:
//
//	len(mockedUnitOfWork.RollbackCalls())
func (mock *UnitOfWorkMock) RollbackCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockRollback.RLock()
	calls = mock.calls.Rollback
	mock.lockRollback.RUnlock()
	return calls
}

// WriteLockBlocksTable calls WriteLockBlocksTableFunc.
func (mock *UnitOfWorkMock) WriteLockBlocksTable(ctx context.Context) error {
	if mock.WriteLockBlocksTableFunc == nil {
		panic("UnitOfWorkMock.WriteLockBlocksTableFunc: method is nil but UnitOfWork.WriteLockBlocksTable was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockWriteLockBlocksTable.Lock()
	mock.calls.WriteLockBlocksTable = append(mock.calls.WriteLockBlocksTable, callInfo)
	mock.lockWriteLockBlocksTable.Unlock()
	return mock.WriteLockBlocksTableFunc(ctx)
}

// WriteLockBlocksTableCalls gets all the calls that were made to WriteLockBlocksTable.
// Check the length with:
//
//	len(mockedUnitOfWork.WriteLockBlocksTableCalls())
func (mock *UnitOfWorkMock) WriteLockBlocksTableCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockWriteLockBlocksTable.RLock()
	calls = mock.calls.WriteLockBlocksTable
	mock.lockWriteLockBlocksTable.RUnlock()
	return calls
}
