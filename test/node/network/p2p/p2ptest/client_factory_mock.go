// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package p2ptest

import (
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"sync"
)

// Ensure, that ClientFactoryMock does implement ClientFactory.
// If this is not the case, regenerate this file with moq.
var _ p2p.ClientFactory = &ClientFactoryMock{}

// ClientFactoryMock is a mock implementation of ClientFactory.
//
// 	func TestSomethingThatUsesClientFactory(t *testing.T) {
//
// 		// make and configure a mocked ClientFactory
// 		mockedClientFactory := &ClientFactoryMock{
// 			CreateClientFunc: func(ip string, port uint16, target string) (Client, error) {
// 				panic("mock out the CreateClient method")
// 			},
// 		}
//
// 		// use mockedClientFactory in code that requires ClientFactory
// 		// and then make assertions.
//
// 	}
type ClientFactoryMock struct {
	// CreateClientFunc mocks the CreateClient method.
	CreateClientFunc func(ip string, port uint16, target string) (p2p.Client, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateClient holds details about calls to the CreateClient method.
		CreateClient []struct {
			// IP is the ip argument value.
			IP string
			// Port is the port argument value.
			Port uint16
			// Target is the target argument value.
			Target string
		}
	}
	lockCreateClient sync.RWMutex
}

// CreateClient calls CreateClientFunc.
func (mock *ClientFactoryMock) CreateClient(ip string, port uint16, target string) (p2p.Client, error) {
	if mock.CreateClientFunc == nil {
		panic("ClientFactoryMock.CreateClientFunc: method is nil but ClientFactory.CreateClient was just called")
	}
	callInfo := struct {
		IP     string
		Port   uint16
		Target string
	}{
		IP:     ip,
		Port:   port,
		Target: target,
	}
	mock.lockCreateClient.Lock()
	mock.calls.CreateClient = append(mock.calls.CreateClient, callInfo)
	mock.lockCreateClient.Unlock()
	return mock.CreateClientFunc(ip, port, target)
}

// CreateClientCalls gets all the calls that were made to CreateClient.
// Check the length with:
//     len(mockedClientFactory.CreateClientCalls())
func (mock *ClientFactoryMock) CreateClientCalls() []struct {
	IP     string
	Port   uint16
	Target string
} {
	var calls []struct {
		IP     string
		Port   uint16
		Target string
	}
	mock.lockCreateClient.RLock()
	calls = mock.calls.CreateClient
	mock.lockCreateClient.RUnlock()
	return calls
}
