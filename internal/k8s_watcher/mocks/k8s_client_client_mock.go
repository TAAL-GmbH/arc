// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/bitcoin-sv/arc/internal/k8s_watcher"
	"sync"
)

// Ensure, that K8sClientMock does implement k8s_watcher.K8sClient.
// If this is not the case, regenerate this file with moq.
var _ k8s_watcher.K8sClient = &K8sClientMock{}

// K8sClientMock is a mock implementation of k8s_watcher.K8sClient.
//
//	func TestSomethingThatUsesK8sClient(t *testing.T) {
//
//		// make and configure a mocked k8s_watcher.K8sClient
//		mockedK8sClient := &K8sClientMock{
//			GetRunningPodNamesFunc: func(ctx context.Context, namespace string, service string) (map[string]struct{}, error) {
//				panic("mock out the GetRunningPodNames method")
//			},
//		}
//
//		// use mockedK8sClient in code that requires k8s_watcher.K8sClient
//		// and then make assertions.
//
//	}
type K8sClientMock struct {
	// GetRunningPodNamesFunc mocks the GetRunningPodNames method.
	GetRunningPodNamesFunc func(ctx context.Context, namespace string, service string) (map[string]struct{}, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetRunningPodNames holds details about calls to the GetRunningPodNames method.
		GetRunningPodNames []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Namespace is the namespace argument value.
			Namespace string
			// Service is the service argument value.
			Service string
		}
	}
	lockGetRunningPodNames sync.RWMutex
}

// GetRunningPodNames calls GetRunningPodNamesFunc.
func (mock *K8sClientMock) GetRunningPodNames(ctx context.Context, namespace string, service string) (map[string]struct{}, error) {
	if mock.GetRunningPodNamesFunc == nil {
		panic("K8sClientMock.GetRunningPodNamesFunc: method is nil but K8sClient.GetRunningPodNames was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Namespace string
		Service   string
	}{
		Ctx:       ctx,
		Namespace: namespace,
		Service:   service,
	}
	mock.lockGetRunningPodNames.Lock()
	mock.calls.GetRunningPodNames = append(mock.calls.GetRunningPodNames, callInfo)
	mock.lockGetRunningPodNames.Unlock()
	return mock.GetRunningPodNamesFunc(ctx, namespace, service)
}

// GetRunningPodNamesCalls gets all the calls that were made to GetRunningPodNames.
// Check the length with:
//
//	len(mockedK8sClient.GetRunningPodNamesCalls())
func (mock *K8sClientMock) GetRunningPodNamesCalls() []struct {
	Ctx       context.Context
	Namespace string
	Service   string
} {
	var calls []struct {
		Ctx       context.Context
		Namespace string
		Service   string
	}
	mock.lockGetRunningPodNames.RLock()
	calls = mock.calls.GetRunningPodNames
	mock.lockGetRunningPodNames.RUnlock()
	return calls
}