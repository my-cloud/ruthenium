// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package p2ptest

import (
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"sync"
)

// Ensure, that ClientMock does implement Client.
// If this is not the case, regenerate this file with moq.
var _ p2p.Client = &ClientMock{}

// ClientMock is a mock implementation of Client.
//
//	func TestSomethingThatUsesClient(t *testing.T) {
//
//		// make and configure a mocked Client
//		mockedClient := &ClientMock{
//			SendFunc: func(topic string, req []byte) ([]byte, error) {
//				panic("mock out the Send method")
//			},
//		}
//
//		// use mockedClient in code that requires Client
//		// and then make assertions.
//
//	}
type ClientMock struct {
	// SendFunc mocks the Send method.
	SendFunc func(topic string, req []byte) ([]byte, error)

	// calls tracks calls to the methods.
	calls struct {
		// Send holds details about calls to the Send method.
		Send []struct {
			// Topic is the topic argument value.
			Topic string
			// Req is the req argument value.
			Req []byte
		}
	}
	lockSend sync.RWMutex
}

// Send calls SendFunc.
func (mock *ClientMock) Send(topic string, req []byte) ([]byte, error) {
	if mock.SendFunc == nil {
		panic("ClientMock.SendFunc: method is nil but Client.Send was just called")
	}
	callInfo := struct {
		Topic string
		Req   []byte
	}{
		Topic: topic,
		Req:   req,
	}
	mock.lockSend.Lock()
	mock.calls.Send = append(mock.calls.Send, callInfo)
	mock.lockSend.Unlock()
	return mock.SendFunc(topic, req)
}

// SendCalls gets all the calls that were made to Send.
// Check the length with:
//
//	len(mockedClient.SendCalls())
func (mock *ClientMock) SendCalls() []struct {
	Topic string
	Req   []byte
} {
	var calls []struct {
		Topic string
		Req   []byte
	}
	mock.lockSend.RLock()
	calls = mock.calls.Send
	mock.lockSend.RUnlock()
	return calls
}
