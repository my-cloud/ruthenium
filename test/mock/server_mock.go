// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	p2p "github.com/leprosus/golang-p2p"
	p2p2 "github.com/my-cloud/ruthenium/src/node/network/p2p"
	"sync"
)

// Ensure, that ServerMock does implement Server.
// If this is not the case, regenerate this file with moq.
var _ p2p2.Server = &ServerMock{}

// ServerMock is a mock implementation of Server.
//
// 	func TestSomethingThatUsesServer(t *testing.T) {
//
// 		// make and configure a mocked Server
// 		mockedServer := &ServerMock{
// 			ServeFunc: func() error {
// 				panic("mock out the Serve method")
// 			},
// 			SetHandleFunc: func(topic string, handler p2p.Handler)  {
// 				panic("mock out the SetHandle method")
// 			},
// 		}
//
// 		// use mockedServer in code that requires Server
// 		// and then make assertions.
//
// 	}
type ServerMock struct {
	// ServeFunc mocks the Serve method.
	ServeFunc func() error

	// SetHandleFunc mocks the SetHandle method.
	SetHandleFunc func(topic string, handler p2p.Handler)

	// calls tracks calls to the methods.
	calls struct {
		// Serve holds details about calls to the Serve method.
		Serve []struct {
		}
		// SetHandle holds details about calls to the SetHandle method.
		SetHandle []struct {
			// Topic is the topic argument value.
			Topic string
			// Handler is the handler argument value.
			Handler p2p.Handler
		}
	}
	lockServe     sync.RWMutex
	lockSetHandle sync.RWMutex
}

// Serve calls ServeFunc.
func (mock *ServerMock) Serve() error {
	if mock.ServeFunc == nil {
		panic("ServerMock.ServeFunc: method is nil but Server.Serve was just called")
	}
	callInfo := struct {
	}{}
	mock.lockServe.Lock()
	mock.calls.Serve = append(mock.calls.Serve, callInfo)
	mock.lockServe.Unlock()
	return mock.ServeFunc()
}

// ServeCalls gets all the calls that were made to Serve.
// Check the length with:
//     len(mockedServer.ServeCalls())
func (mock *ServerMock) ServeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockServe.RLock()
	calls = mock.calls.Serve
	mock.lockServe.RUnlock()
	return calls
}

// SetHandle calls SetHandleFunc.
func (mock *ServerMock) SetHandle(topic string, handler p2p.Handler) {
	if mock.SetHandleFunc == nil {
		panic("ServerMock.SetHandleFunc: method is nil but Server.SetHandle was just called")
	}
	callInfo := struct {
		Topic   string
		Handler p2p.Handler
	}{
		Topic:   topic,
		Handler: handler,
	}
	mock.lockSetHandle.Lock()
	mock.calls.SetHandle = append(mock.calls.SetHandle, callInfo)
	mock.lockSetHandle.Unlock()
	mock.SetHandleFunc(topic, handler)
}

// SetHandleCalls gets all the calls that were made to SetHandle.
// Check the length with:
//     len(mockedServer.SetHandleCalls())
func (mock *ServerMock) SetHandleCalls() []struct {
	Topic   string
	Handler p2p.Handler
} {
	var calls []struct {
		Topic   string
		Handler p2p.Handler
	}
	mock.lockSetHandle.RLock()
	calls = mock.calls.SetHandle
	mock.lockSetHandle.RUnlock()
	return calls
}
