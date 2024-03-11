// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package blocktx

import (
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/libsv/go-p2p/wire"
	"sync"
)

// Ensure, that PeerMock does implement Peer.
// If this is not the case, regenerate this file with moq.
var _ Peer = &PeerMock{}

// PeerMock is a mock implementation of Peer.
//
//	func TestSomethingThatUsesPeer(t *testing.T) {
//
//		// make and configure a mocked Peer
//		mockedPeer := &PeerMock{
//			AnnounceBlockFunc: func(blockHash *chainhash.Hash)  {
//				panic("mock out the AnnounceBlock method")
//			},
//			AnnounceTransactionFunc: func(txHash *chainhash.Hash)  {
//				panic("mock out the AnnounceTransaction method")
//			},
//			ConnectedFunc: func() bool {
//				panic("mock out the Connected method")
//			},
//			IsHealthyFunc: func() bool {
//				panic("mock out the IsHealthy method")
//			},
//			NetworkFunc: func() wire.BitcoinNet {
//				panic("mock out the Network method")
//			},
//			RequestBlockFunc: func(blockHash *chainhash.Hash)  {
//				panic("mock out the RequestBlock method")
//			},
//			RequestTransactionFunc: func(txHash *chainhash.Hash)  {
//				panic("mock out the RequestTransaction method")
//			},
//			StringFunc: func() string {
//				panic("mock out the String method")
//			},
//			WriteMsgFunc: func(msg wire.Message) error {
//				panic("mock out the WriteMsg method")
//			},
//		}
//
//		// use mockedPeer in code that requires Peer
//		// and then make assertions.
//
//	}
type PeerMock struct {
	// AnnounceBlockFunc mocks the AnnounceBlock method.
	AnnounceBlockFunc func(blockHash *chainhash.Hash)

	// AnnounceTransactionFunc mocks the AnnounceTransaction method.
	AnnounceTransactionFunc func(txHash *chainhash.Hash)

	// ConnectedFunc mocks the Connected method.
	ConnectedFunc func() bool

	// IsHealthyFunc mocks the IsHealthy method.
	IsHealthyFunc func() bool

	// NetworkFunc mocks the Network method.
	NetworkFunc func() wire.BitcoinNet

	// RequestBlockFunc mocks the RequestBlock method.
	RequestBlockFunc func(blockHash *chainhash.Hash)

	// RequestTransactionFunc mocks the RequestTransaction method.
	RequestTransactionFunc func(txHash *chainhash.Hash)

	// StringFunc mocks the String method.
	StringFunc func() string

	// WriteMsgFunc mocks the WriteMsg method.
	WriteMsgFunc func(msg wire.Message) error

	// calls tracks calls to the methods.
	calls struct {
		// AnnounceBlock holds details about calls to the AnnounceBlock method.
		AnnounceBlock []struct {
			// BlockHash is the blockHash argument value.
			BlockHash *chainhash.Hash
		}
		// AnnounceTransaction holds details about calls to the AnnounceTransaction method.
		AnnounceTransaction []struct {
			// TxHash is the txHash argument value.
			TxHash *chainhash.Hash
		}
		// Connected holds details about calls to the Connected method.
		Connected []struct {
		}
		// IsHealthy holds details about calls to the IsHealthy method.
		IsHealthy []struct {
		}
		// Network holds details about calls to the Network method.
		Network []struct {
		}
		// RequestBlock holds details about calls to the RequestBlock method.
		RequestBlock []struct {
			// BlockHash is the blockHash argument value.
			BlockHash *chainhash.Hash
		}
		// RequestTransaction holds details about calls to the RequestTransaction method.
		RequestTransaction []struct {
			// TxHash is the txHash argument value.
			TxHash *chainhash.Hash
		}
		// String holds details about calls to the String method.
		String []struct {
		}
		// WriteMsg holds details about calls to the WriteMsg method.
		WriteMsg []struct {
			// Msg is the msg argument value.
			Msg wire.Message
		}
	}
	lockAnnounceBlock       sync.RWMutex
	lockAnnounceTransaction sync.RWMutex
	lockConnected           sync.RWMutex
	lockIsHealthy           sync.RWMutex
	lockNetwork             sync.RWMutex
	lockRequestBlock        sync.RWMutex
	lockRequestTransaction  sync.RWMutex
	lockString              sync.RWMutex
	lockWriteMsg            sync.RWMutex
}

// AnnounceBlock calls AnnounceBlockFunc.
func (mock *PeerMock) AnnounceBlock(blockHash *chainhash.Hash) {
	if mock.AnnounceBlockFunc == nil {
		panic("PeerMock.AnnounceBlockFunc: method is nil but Peer.AnnounceBlock was just called")
	}
	callInfo := struct {
		BlockHash *chainhash.Hash
	}{
		BlockHash: blockHash,
	}
	mock.lockAnnounceBlock.Lock()
	mock.calls.AnnounceBlock = append(mock.calls.AnnounceBlock, callInfo)
	mock.lockAnnounceBlock.Unlock()
	mock.AnnounceBlockFunc(blockHash)
}

// AnnounceBlockCalls gets all the calls that were made to AnnounceBlock.
// Check the length with:
//
//	len(mockedPeer.AnnounceBlockCalls())
func (mock *PeerMock) AnnounceBlockCalls() []struct {
	BlockHash *chainhash.Hash
} {
	var calls []struct {
		BlockHash *chainhash.Hash
	}
	mock.lockAnnounceBlock.RLock()
	calls = mock.calls.AnnounceBlock
	mock.lockAnnounceBlock.RUnlock()
	return calls
}

// AnnounceTransaction calls AnnounceTransactionFunc.
func (mock *PeerMock) AnnounceTransaction(txHash *chainhash.Hash) {
	if mock.AnnounceTransactionFunc == nil {
		panic("PeerMock.AnnounceTransactionFunc: method is nil but Peer.AnnounceTransaction was just called")
	}
	callInfo := struct {
		TxHash *chainhash.Hash
	}{
		TxHash: txHash,
	}
	mock.lockAnnounceTransaction.Lock()
	mock.calls.AnnounceTransaction = append(mock.calls.AnnounceTransaction, callInfo)
	mock.lockAnnounceTransaction.Unlock()
	mock.AnnounceTransactionFunc(txHash)
}

// AnnounceTransactionCalls gets all the calls that were made to AnnounceTransaction.
// Check the length with:
//
//	len(mockedPeer.AnnounceTransactionCalls())
func (mock *PeerMock) AnnounceTransactionCalls() []struct {
	TxHash *chainhash.Hash
} {
	var calls []struct {
		TxHash *chainhash.Hash
	}
	mock.lockAnnounceTransaction.RLock()
	calls = mock.calls.AnnounceTransaction
	mock.lockAnnounceTransaction.RUnlock()
	return calls
}

// Connected calls ConnectedFunc.
func (mock *PeerMock) Connected() bool {
	if mock.ConnectedFunc == nil {
		panic("PeerMock.ConnectedFunc: method is nil but Peer.Connected was just called")
	}
	callInfo := struct {
	}{}
	mock.lockConnected.Lock()
	mock.calls.Connected = append(mock.calls.Connected, callInfo)
	mock.lockConnected.Unlock()
	return mock.ConnectedFunc()
}

// ConnectedCalls gets all the calls that were made to Connected.
// Check the length with:
//
//	len(mockedPeer.ConnectedCalls())
func (mock *PeerMock) ConnectedCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockConnected.RLock()
	calls = mock.calls.Connected
	mock.lockConnected.RUnlock()
	return calls
}

// IsHealthy calls IsHealthyFunc.
func (mock *PeerMock) IsHealthy() bool {
	if mock.IsHealthyFunc == nil {
		panic("PeerMock.IsHealthyFunc: method is nil but Peer.IsHealthy was just called")
	}
	callInfo := struct {
	}{}
	mock.lockIsHealthy.Lock()
	mock.calls.IsHealthy = append(mock.calls.IsHealthy, callInfo)
	mock.lockIsHealthy.Unlock()
	return mock.IsHealthyFunc()
}

// IsHealthyCalls gets all the calls that were made to IsHealthy.
// Check the length with:
//
//	len(mockedPeer.IsHealthyCalls())
func (mock *PeerMock) IsHealthyCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIsHealthy.RLock()
	calls = mock.calls.IsHealthy
	mock.lockIsHealthy.RUnlock()
	return calls
}

// Network calls NetworkFunc.
func (mock *PeerMock) Network() wire.BitcoinNet {
	if mock.NetworkFunc == nil {
		panic("PeerMock.NetworkFunc: method is nil but Peer.Network was just called")
	}
	callInfo := struct {
	}{}
	mock.lockNetwork.Lock()
	mock.calls.Network = append(mock.calls.Network, callInfo)
	mock.lockNetwork.Unlock()
	return mock.NetworkFunc()
}

// NetworkCalls gets all the calls that were made to Network.
// Check the length with:
//
//	len(mockedPeer.NetworkCalls())
func (mock *PeerMock) NetworkCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockNetwork.RLock()
	calls = mock.calls.Network
	mock.lockNetwork.RUnlock()
	return calls
}

// RequestBlock calls RequestBlockFunc.
func (mock *PeerMock) RequestBlock(blockHash *chainhash.Hash) {
	if mock.RequestBlockFunc == nil {
		panic("PeerMock.RequestBlockFunc: method is nil but Peer.RequestBlock was just called")
	}
	callInfo := struct {
		BlockHash *chainhash.Hash
	}{
		BlockHash: blockHash,
	}
	mock.lockRequestBlock.Lock()
	mock.calls.RequestBlock = append(mock.calls.RequestBlock, callInfo)
	mock.lockRequestBlock.Unlock()
	mock.RequestBlockFunc(blockHash)
}

// RequestBlockCalls gets all the calls that were made to RequestBlock.
// Check the length with:
//
//	len(mockedPeer.RequestBlockCalls())
func (mock *PeerMock) RequestBlockCalls() []struct {
	BlockHash *chainhash.Hash
} {
	var calls []struct {
		BlockHash *chainhash.Hash
	}
	mock.lockRequestBlock.RLock()
	calls = mock.calls.RequestBlock
	mock.lockRequestBlock.RUnlock()
	return calls
}

// RequestTransaction calls RequestTransactionFunc.
func (mock *PeerMock) RequestTransaction(txHash *chainhash.Hash) {
	if mock.RequestTransactionFunc == nil {
		panic("PeerMock.RequestTransactionFunc: method is nil but Peer.RequestTransaction was just called")
	}
	callInfo := struct {
		TxHash *chainhash.Hash
	}{
		TxHash: txHash,
	}
	mock.lockRequestTransaction.Lock()
	mock.calls.RequestTransaction = append(mock.calls.RequestTransaction, callInfo)
	mock.lockRequestTransaction.Unlock()
	mock.RequestTransactionFunc(txHash)
}

// RequestTransactionCalls gets all the calls that were made to RequestTransaction.
// Check the length with:
//
//	len(mockedPeer.RequestTransactionCalls())
func (mock *PeerMock) RequestTransactionCalls() []struct {
	TxHash *chainhash.Hash
} {
	var calls []struct {
		TxHash *chainhash.Hash
	}
	mock.lockRequestTransaction.RLock()
	calls = mock.calls.RequestTransaction
	mock.lockRequestTransaction.RUnlock()
	return calls
}

// String calls StringFunc.
func (mock *PeerMock) String() string {
	if mock.StringFunc == nil {
		panic("PeerMock.StringFunc: method is nil but Peer.String was just called")
	}
	callInfo := struct {
	}{}
	mock.lockString.Lock()
	mock.calls.String = append(mock.calls.String, callInfo)
	mock.lockString.Unlock()
	return mock.StringFunc()
}

// StringCalls gets all the calls that were made to String.
// Check the length with:
//
//	len(mockedPeer.StringCalls())
func (mock *PeerMock) StringCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockString.RLock()
	calls = mock.calls.String
	mock.lockString.RUnlock()
	return calls
}

// WriteMsg calls WriteMsgFunc.
func (mock *PeerMock) WriteMsg(msg wire.Message) error {
	if mock.WriteMsgFunc == nil {
		panic("PeerMock.WriteMsgFunc: method is nil but Peer.WriteMsg was just called")
	}
	callInfo := struct {
		Msg wire.Message
	}{
		Msg: msg,
	}
	mock.lockWriteMsg.Lock()
	mock.calls.WriteMsg = append(mock.calls.WriteMsg, callInfo)
	mock.lockWriteMsg.Unlock()
	return mock.WriteMsgFunc(msg)
}

// WriteMsgCalls gets all the calls that were made to WriteMsg.
// Check the length with:
//
//	len(mockedPeer.WriteMsgCalls())
func (mock *PeerMock) WriteMsgCalls() []struct {
	Msg wire.Message
} {
	var calls []struct {
		Msg wire.Message
	}
	mock.lockWriteMsg.RLock()
	calls = mock.calls.WriteMsg
	mock.lockWriteMsg.RUnlock()
	return calls
}
