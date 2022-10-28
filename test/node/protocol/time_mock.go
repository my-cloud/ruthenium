// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/clock"
	"sync"
	"time"
)

// Ensure, that TimeMock does implement Time.
// If this is not the case, regenerate this file with moq.
var _ clock.Time = &TimeMock{}

// TimeMock is a mock implementation of Time.
//
// 	func TestSomethingThatUsesTime(t *testing.T) {
//
// 		// make and configure a mocked Time
// 		mockedTime := &TimeMock{
// 			NowFunc: func() time.Time {
// 				panic("mock out the Now method")
// 			},
// 		}
//
// 		// use mockedTime in code that requires Time
// 		// and then make assertions.
//
// 	}
type TimeMock struct {
	// NowFunc mocks the Now method.
	NowFunc func() time.Time

	// calls tracks calls to the methods.
	calls struct {
		// Now holds details about calls to the Now method.
		Now []struct {
		}
	}
	lockNow sync.RWMutex
}

// Now calls NowFunc.
func (mock *TimeMock) Now() time.Time {
	if mock.NowFunc == nil {
		panic("TimeMock.NowFunc: method is nil but Time.Now was just called")
	}
	callInfo := struct {
	}{}
	mock.lockNow.Lock()
	mock.calls.Now = append(mock.calls.Now, callInfo)
	mock.lockNow.Unlock()
	return mock.NowFunc()
}

// NowCalls gets all the calls that were made to Now.
// Check the length with:
//     len(mockedTime.NowCalls())
func (mock *TimeMock) NowCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockNow.RLock()
	calls = mock.calls.Now
	mock.lockNow.RUnlock()
	return calls
}
