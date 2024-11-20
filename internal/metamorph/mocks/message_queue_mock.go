// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/bitcoin-sv/arc/internal/metamorph"
	"sync"
)

// Ensure, that MessageQueueClientMock does implement metamorph.MessageQueueClient.
// If this is not the case, regenerate this file with moq.
var _ metamorph.MessageQueueClient = &MessageQueueClientMock{}

// MessageQueueClientMock is a mock implementation of metamorph.MessageQueueClient.
//
//	func TestSomethingThatUsesMessageQueueClient(t *testing.T) {
//
//		// make and configure a mocked metamorph.MessageQueueClient
//		mockedMessageQueueClient := &MessageQueueClientMock{
//			PublishFunc: func(ctx context.Context, topic string, data []byte) error {
//				panic("mock out the Publish method")
//			},
//			ShutdownFunc: func()  {
//				panic("mock out the Shutdown method")
//			},
//			SubscribeFunc: func(topic string, msgFunc func([]byte) error) error {
//				panic("mock out the Subscribe method")
//			},
//		}
//
//		// use mockedMessageQueueClient in code that requires metamorph.MessageQueueClient
//		// and then make assertions.
//
//	}
type MessageQueueClientMock struct {
	// PublishFunc mocks the Publish method.
	PublishFunc func(ctx context.Context, topic string, data []byte) error

	// ShutdownFunc mocks the Shutdown method.
	ShutdownFunc func()

	// SubscribeFunc mocks the Subscribe method.
	SubscribeFunc func(topic string, msgFunc func([]byte) error) error

	// calls tracks calls to the methods.
	calls struct {
		// Publish holds details about calls to the Publish method.
		Publish []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Topic is the topic argument value.
			Topic string
			// Data is the data argument value.
			Data []byte
		}
		// Shutdown holds details about calls to the Shutdown method.
		Shutdown []struct {
		}
		// Subscribe holds details about calls to the Subscribe method.
		Subscribe []struct {
			// Topic is the topic argument value.
			Topic string
			// MsgFunc is the msgFunc argument value.
			MsgFunc func([]byte) error
		}
	}
	lockPublish   sync.RWMutex
	lockShutdown  sync.RWMutex
	lockSubscribe sync.RWMutex
}

// Publish calls PublishFunc.
func (mock *MessageQueueClientMock) Publish(ctx context.Context, topic string, data []byte) error {
	if mock.PublishFunc == nil {
		panic("MessageQueueClientMock.PublishFunc: method is nil but MessageQueueClient.Publish was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Topic string
		Data  []byte
	}{
		Ctx:   ctx,
		Topic: topic,
		Data:  data,
	}
	mock.lockPublish.Lock()
	mock.calls.Publish = append(mock.calls.Publish, callInfo)
	mock.lockPublish.Unlock()
	return mock.PublishFunc(ctx, topic, data)
}

// PublishCalls gets all the calls that were made to Publish.
// Check the length with:
//
//	len(mockedMessageQueueClient.PublishCalls())
func (mock *MessageQueueClientMock) PublishCalls() []struct {
	Ctx   context.Context
	Topic string
	Data  []byte
} {
	var calls []struct {
		Ctx   context.Context
		Topic string
		Data  []byte
	}
	mock.lockPublish.RLock()
	calls = mock.calls.Publish
	mock.lockPublish.RUnlock()
	return calls
}

// Shutdown calls ShutdownFunc.
func (mock *MessageQueueClientMock) Shutdown() {
	if mock.ShutdownFunc == nil {
		panic("MessageQueueClientMock.ShutdownFunc: method is nil but MessageQueueClient.Shutdown was just called")
	}
	callInfo := struct {
	}{}
	mock.lockShutdown.Lock()
	mock.calls.Shutdown = append(mock.calls.Shutdown, callInfo)
	mock.lockShutdown.Unlock()
	mock.ShutdownFunc()
}

// ShutdownCalls gets all the calls that were made to Shutdown.
// Check the length with:
//
//	len(mockedMessageQueueClient.ShutdownCalls())
func (mock *MessageQueueClientMock) ShutdownCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockShutdown.RLock()
	calls = mock.calls.Shutdown
	mock.lockShutdown.RUnlock()
	return calls
}

// Subscribe calls SubscribeFunc.
func (mock *MessageQueueClientMock) Subscribe(topic string, msgFunc func([]byte) error) error {
	if mock.SubscribeFunc == nil {
		panic("MessageQueueClientMock.SubscribeFunc: method is nil but MessageQueueClient.Subscribe was just called")
	}
	callInfo := struct {
		Topic   string
		MsgFunc func([]byte) error
	}{
		Topic:   topic,
		MsgFunc: msgFunc,
	}
	mock.lockSubscribe.Lock()
	mock.calls.Subscribe = append(mock.calls.Subscribe, callInfo)
	mock.lockSubscribe.Unlock()
	return mock.SubscribeFunc(topic, msgFunc)
}

// SubscribeCalls gets all the calls that were made to Subscribe.
// Check the length with:
//
//	len(mockedMessageQueueClient.SubscribeCalls())
func (mock *MessageQueueClientMock) SubscribeCalls() []struct {
	Topic   string
	MsgFunc func([]byte) error
} {
	var calls []struct {
		Topic   string
		MsgFunc func([]byte) error
	}
	mock.lockSubscribe.RLock()
	calls = mock.calls.Subscribe
	mock.lockSubscribe.RUnlock()
	return calls
}
