// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package networktest

import (
	"github.com/my-cloud/ruthenium/src/node/network"
	"sync"
)

// Ensure, that IpFinderMock does implement IpFinder.
// If this is not the case, regenerate this file with moq.
var _ network.IpFinder = &IpFinderMock{}

// IpFinderMock is a mock implementation of IpFinder.
//
//	func TestSomethingThatUsesIpFinder(t *testing.T) {
//
//		// make and configure a mocked IpFinder
//		mockedIpFinder := &IpFinderMock{
//			LookupIPFunc: func(ip string) (string, error) {
//				panic("mock out the LookupIP method")
//			},
//		}
//
//		// use mockedIpFinder in code that requires IpFinder
//		// and then make assertions.
//
//	}
type IpFinderMock struct {
	// LookupIPFunc mocks the LookupIP method.
	LookupIPFunc func(ip string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// LookupIP holds details about calls to the LookupIP method.
		LookupIP []struct {
			// IP is the ip argument value.
			IP string
		}
	}
	lockLookupIP sync.RWMutex
}

// LookupIP calls LookupIPFunc.
func (mock *IpFinderMock) LookupIP(ip string) (string, error) {
	if mock.LookupIPFunc == nil {
		panic("IpFinderMock.LookupIPFunc: method is nil but IpFinder.LookupIP was just called")
	}
	callInfo := struct {
		IP string
	}{
		IP: ip,
	}
	mock.lockLookupIP.Lock()
	mock.calls.LookupIP = append(mock.calls.LookupIP, callInfo)
	mock.lockLookupIP.Unlock()
	return mock.LookupIPFunc(ip)
}

// LookupIPCalls gets all the calls that were made to LookupIP.
// Check the length with:
//
//	len(mockedIpFinder.LookupIPCalls())
func (mock *IpFinderMock) LookupIPCalls() []struct {
	IP string
} {
	var calls []struct {
		IP string
	}
	mock.lockLookupIP.RLock()
	calls = mock.calls.LookupIP
	mock.lockLookupIP.RUnlock()
	return calls
}
