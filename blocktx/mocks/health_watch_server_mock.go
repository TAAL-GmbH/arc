// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/bitcoin-sv/arc/blocktx"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"sync"
)

// Ensure, that HealthWatchServerMock does implement blocktx.HealthWatchServer.
// If this is not the case, regenerate this file with moq.
var _ blocktx.HealthWatchServer = &HealthWatchServerMock{}

// HealthWatchServerMock is a mock implementation of blocktx.HealthWatchServer.
//
//	func TestSomethingThatUsesHealthWatchServer(t *testing.T) {
//
//		// make and configure a mocked blocktx.HealthWatchServer
//		mockedHealthWatchServer := &HealthWatchServerMock{
//			ContextFunc: func() context.Context {
//				panic("mock out the Context method")
//			},
//			RecvMsgFunc: func(m interface{}) error {
//				panic("mock out the RecvMsg method")
//			},
//			SendFunc: func(healthCheckResponse *grpc_health_v1.HealthCheckResponse) error {
//				panic("mock out the Send method")
//			},
//			SendHeaderFunc: func(mD metadata.MD) error {
//				panic("mock out the SendHeader method")
//			},
//			SendMsgFunc: func(m interface{}) error {
//				panic("mock out the SendMsg method")
//			},
//			SetHeaderFunc: func(mD metadata.MD) error {
//				panic("mock out the SetHeader method")
//			},
//			SetTrailerFunc: func(mD metadata.MD)  {
//				panic("mock out the SetTrailer method")
//			},
//		}
//
//		// use mockedHealthWatchServer in code that requires blocktx.HealthWatchServer
//		// and then make assertions.
//
//	}
type HealthWatchServerMock struct {
	// ContextFunc mocks the Context method.
	ContextFunc func() context.Context

	// RecvMsgFunc mocks the RecvMsg method.
	RecvMsgFunc func(m interface{}) error

	// SendFunc mocks the Send method.
	SendFunc func(healthCheckResponse *grpc_health_v1.HealthCheckResponse) error

	// SendHeaderFunc mocks the SendHeader method.
	SendHeaderFunc func(mD metadata.MD) error

	// SendMsgFunc mocks the SendMsg method.
	SendMsgFunc func(m interface{}) error

	// SetHeaderFunc mocks the SetHeader method.
	SetHeaderFunc func(mD metadata.MD) error

	// SetTrailerFunc mocks the SetTrailer method.
	SetTrailerFunc func(mD metadata.MD)

	// calls tracks calls to the methods.
	calls struct {
		// Context holds details about calls to the Context method.
		Context []struct {
		}
		// RecvMsg holds details about calls to the RecvMsg method.
		RecvMsg []struct {
			// M is the m argument value.
			M interface{}
		}
		// Send holds details about calls to the Send method.
		Send []struct {
			// HealthCheckResponse is the healthCheckResponse argument value.
			HealthCheckResponse *grpc_health_v1.HealthCheckResponse
		}
		// SendHeader holds details about calls to the SendHeader method.
		SendHeader []struct {
			// MD is the mD argument value.
			MD metadata.MD
		}
		// SendMsg holds details about calls to the SendMsg method.
		SendMsg []struct {
			// M is the m argument value.
			M interface{}
		}
		// SetHeader holds details about calls to the SetHeader method.
		SetHeader []struct {
			// MD is the mD argument value.
			MD metadata.MD
		}
		// SetTrailer holds details about calls to the SetTrailer method.
		SetTrailer []struct {
			// MD is the mD argument value.
			MD metadata.MD
		}
	}
	lockContext    sync.RWMutex
	lockRecvMsg    sync.RWMutex
	lockSend       sync.RWMutex
	lockSendHeader sync.RWMutex
	lockSendMsg    sync.RWMutex
	lockSetHeader  sync.RWMutex
	lockSetTrailer sync.RWMutex
}

// Context calls ContextFunc.
func (mock *HealthWatchServerMock) Context() context.Context {
	if mock.ContextFunc == nil {
		panic("HealthWatchServerMock.ContextFunc: method is nil but HealthWatchServer.Context was just called")
	}
	callInfo := struct {
	}{}
	mock.lockContext.Lock()
	mock.calls.Context = append(mock.calls.Context, callInfo)
	mock.lockContext.Unlock()
	return mock.ContextFunc()
}

// ContextCalls gets all the calls that were made to Context.
// Check the length with:
//
//	len(mockedHealthWatchServer.ContextCalls())
func (mock *HealthWatchServerMock) ContextCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockContext.RLock()
	calls = mock.calls.Context
	mock.lockContext.RUnlock()
	return calls
}

// RecvMsg calls RecvMsgFunc.
func (mock *HealthWatchServerMock) RecvMsg(m interface{}) error {
	if mock.RecvMsgFunc == nil {
		panic("HealthWatchServerMock.RecvMsgFunc: method is nil but HealthWatchServer.RecvMsg was just called")
	}
	callInfo := struct {
		M interface{}
	}{
		M: m,
	}
	mock.lockRecvMsg.Lock()
	mock.calls.RecvMsg = append(mock.calls.RecvMsg, callInfo)
	mock.lockRecvMsg.Unlock()
	return mock.RecvMsgFunc(m)
}

// RecvMsgCalls gets all the calls that were made to RecvMsg.
// Check the length with:
//
//	len(mockedHealthWatchServer.RecvMsgCalls())
func (mock *HealthWatchServerMock) RecvMsgCalls() []struct {
	M interface{}
} {
	var calls []struct {
		M interface{}
	}
	mock.lockRecvMsg.RLock()
	calls = mock.calls.RecvMsg
	mock.lockRecvMsg.RUnlock()
	return calls
}

// Send calls SendFunc.
func (mock *HealthWatchServerMock) Send(healthCheckResponse *grpc_health_v1.HealthCheckResponse) error {
	if mock.SendFunc == nil {
		panic("HealthWatchServerMock.SendFunc: method is nil but HealthWatchServer.Send was just called")
	}
	callInfo := struct {
		HealthCheckResponse *grpc_health_v1.HealthCheckResponse
	}{
		HealthCheckResponse: healthCheckResponse,
	}
	mock.lockSend.Lock()
	mock.calls.Send = append(mock.calls.Send, callInfo)
	mock.lockSend.Unlock()
	return mock.SendFunc(healthCheckResponse)
}

// SendCalls gets all the calls that were made to Send.
// Check the length with:
//
//	len(mockedHealthWatchServer.SendCalls())
func (mock *HealthWatchServerMock) SendCalls() []struct {
	HealthCheckResponse *grpc_health_v1.HealthCheckResponse
} {
	var calls []struct {
		HealthCheckResponse *grpc_health_v1.HealthCheckResponse
	}
	mock.lockSend.RLock()
	calls = mock.calls.Send
	mock.lockSend.RUnlock()
	return calls
}

// SendHeader calls SendHeaderFunc.
func (mock *HealthWatchServerMock) SendHeader(mD metadata.MD) error {
	if mock.SendHeaderFunc == nil {
		panic("HealthWatchServerMock.SendHeaderFunc: method is nil but HealthWatchServer.SendHeader was just called")
	}
	callInfo := struct {
		MD metadata.MD
	}{
		MD: mD,
	}
	mock.lockSendHeader.Lock()
	mock.calls.SendHeader = append(mock.calls.SendHeader, callInfo)
	mock.lockSendHeader.Unlock()
	return mock.SendHeaderFunc(mD)
}

// SendHeaderCalls gets all the calls that were made to SendHeader.
// Check the length with:
//
//	len(mockedHealthWatchServer.SendHeaderCalls())
func (mock *HealthWatchServerMock) SendHeaderCalls() []struct {
	MD metadata.MD
} {
	var calls []struct {
		MD metadata.MD
	}
	mock.lockSendHeader.RLock()
	calls = mock.calls.SendHeader
	mock.lockSendHeader.RUnlock()
	return calls
}

// SendMsg calls SendMsgFunc.
func (mock *HealthWatchServerMock) SendMsg(m interface{}) error {
	if mock.SendMsgFunc == nil {
		panic("HealthWatchServerMock.SendMsgFunc: method is nil but HealthWatchServer.SendMsg was just called")
	}
	callInfo := struct {
		M interface{}
	}{
		M: m,
	}
	mock.lockSendMsg.Lock()
	mock.calls.SendMsg = append(mock.calls.SendMsg, callInfo)
	mock.lockSendMsg.Unlock()
	return mock.SendMsgFunc(m)
}

// SendMsgCalls gets all the calls that were made to SendMsg.
// Check the length with:
//
//	len(mockedHealthWatchServer.SendMsgCalls())
func (mock *HealthWatchServerMock) SendMsgCalls() []struct {
	M interface{}
} {
	var calls []struct {
		M interface{}
	}
	mock.lockSendMsg.RLock()
	calls = mock.calls.SendMsg
	mock.lockSendMsg.RUnlock()
	return calls
}

// SetHeader calls SetHeaderFunc.
func (mock *HealthWatchServerMock) SetHeader(mD metadata.MD) error {
	if mock.SetHeaderFunc == nil {
		panic("HealthWatchServerMock.SetHeaderFunc: method is nil but HealthWatchServer.SetHeader was just called")
	}
	callInfo := struct {
		MD metadata.MD
	}{
		MD: mD,
	}
	mock.lockSetHeader.Lock()
	mock.calls.SetHeader = append(mock.calls.SetHeader, callInfo)
	mock.lockSetHeader.Unlock()
	return mock.SetHeaderFunc(mD)
}

// SetHeaderCalls gets all the calls that were made to SetHeader.
// Check the length with:
//
//	len(mockedHealthWatchServer.SetHeaderCalls())
func (mock *HealthWatchServerMock) SetHeaderCalls() []struct {
	MD metadata.MD
} {
	var calls []struct {
		MD metadata.MD
	}
	mock.lockSetHeader.RLock()
	calls = mock.calls.SetHeader
	mock.lockSetHeader.RUnlock()
	return calls
}

// SetTrailer calls SetTrailerFunc.
func (mock *HealthWatchServerMock) SetTrailer(mD metadata.MD) {
	if mock.SetTrailerFunc == nil {
		panic("HealthWatchServerMock.SetTrailerFunc: method is nil but HealthWatchServer.SetTrailer was just called")
	}
	callInfo := struct {
		MD metadata.MD
	}{
		MD: mD,
	}
	mock.lockSetTrailer.Lock()
	mock.calls.SetTrailer = append(mock.calls.SetTrailer, callInfo)
	mock.lockSetTrailer.Unlock()
	mock.SetTrailerFunc(mD)
}

// SetTrailerCalls gets all the calls that were made to SetTrailer.
// Check the length with:
//
//	len(mockedHealthWatchServer.SetTrailerCalls())
func (mock *HealthWatchServerMock) SetTrailerCalls() []struct {
	MD metadata.MD
} {
	var calls []struct {
		MD metadata.MD
	}
	mock.lockSetTrailer.RLock()
	calls = mock.calls.SetTrailer
	mock.lockSetTrailer.RUnlock()
	return calls
}
