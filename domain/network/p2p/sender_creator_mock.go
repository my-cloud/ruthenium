// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package p2p

import (
	"sync"
)

// Ensure, that SenderCreatorMock does implement SenderCreator.
// If this is not the case, regenerate this file with moq.
var _ SenderCreator = &SenderCreatorMock{}

// SenderCreatorMock is a mock implementation of SenderCreator.
//
//	func TestSomethingThatUsesSenderCreator(t *testing.T) {
//
//		// make and configure a mocked SenderCreator
//		mockedSenderCreator := &SenderCreatorMock{
//			CreateSenderFunc: func(ip string, port string) (Sender, error) {
//				panic("mock out the CreateSender method")
//			},
//		}
//
//		// use mockedSenderCreator in code that requires SenderCreator
//		// and then make assertions.
//
//	}
type SenderCreatorMock struct {
	// CreateSenderFunc mocks the CreateSender method.
	CreateSenderFunc func(ip string, port string) (Sender, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateSender holds details about calls to the CreateSender method.
		CreateSender []struct {
			// IP is the ip argument value.
			IP string
			// Port is the port argument value.
			Port string
		}
	}
	lockCreateSender sync.RWMutex
}

// CreateSender calls CreateSenderFunc.
func (mock *SenderCreatorMock) CreateSender(ip string, port string) (Sender, error) {
	if mock.CreateSenderFunc == nil {
		panic("SenderCreatorMock.CreateSenderFunc: method is nil but SenderCreator.CreateSender was just called")
	}
	callInfo := struct {
		IP   string
		Port string
	}{
		IP:   ip,
		Port: port,
	}
	mock.lockCreateSender.Lock()
	mock.calls.CreateSender = append(mock.calls.CreateSender, callInfo)
	mock.lockCreateSender.Unlock()
	return mock.CreateSenderFunc(ip, port)
}

// CreateSenderCalls gets all the calls that were made to CreateSender.
// Check the length with:
//
//	len(mockedSenderCreator.CreateSenderCalls())
func (mock *SenderCreatorMock) CreateSenderCalls() []struct {
	IP   string
	Port string
} {
	var calls []struct {
		IP   string
		Port string
	}
	mock.lockCreateSender.RLock()
	calls = mock.calls.CreateSender
	mock.lockCreateSender.RUnlock()
	return calls
}
