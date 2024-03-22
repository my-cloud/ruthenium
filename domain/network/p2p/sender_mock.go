// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package p2p

import (
	"sync"
)

// Ensure, that SenderMock does implement Sender.
// If this is not the case, regenerate this file with moq.
var _ Sender = &SenderMock{}

// SenderMock is a mock implementation of Sender.
//
//	func TestSomethingThatUsesSender(t *testing.T) {
//
//		// make and configure a mocked Sender
//		mockedSender := &SenderMock{
//			SendFunc: func(topic string, req []byte) ([]byte, error) {
//				panic("mock out the Send method")
//			},
//		}
//
//		// use mockedSender in code that requires Sender
//		// and then make assertions.
//
//	}
type SenderMock struct {
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
func (mock *SenderMock) Send(topic string, req []byte) ([]byte, error) {
	if mock.SendFunc == nil {
		panic("SenderMock.SendFunc: method is nil but Sender.Send was just called")
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
//	len(mockedSender.SendCalls())
func (mock *SenderMock) SendCalls() []struct {
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