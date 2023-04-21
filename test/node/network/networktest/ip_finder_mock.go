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
//			FindHostPublicIpFunc: func() (string, error) {
//				panic("mock out the FindHostPublicIp method")
//			},
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
	// FindHostPublicIpFunc mocks the FindHostPublicIp method.
	FindHostPublicIpFunc func() (string, error)

	// LookupIPFunc mocks the LookupIP method.
	LookupIPFunc func(ip string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// FindHostPublicIp holds details about calls to the FindHostPublicIp method.
		FindHostPublicIp []struct {
		}
		// LookupIP holds details about calls to the LookupIP method.
		LookupIP []struct {
			// IP is the ip argument value.
			IP string
		}
	}
	lockFindHostPublicIp sync.RWMutex
	lockLookupIP         sync.RWMutex
}

// FindHostPublicIp calls FindHostPublicIpFunc.
func (mock *IpFinderMock) FindHostPublicIp() (string, error) {
	if mock.FindHostPublicIpFunc == nil {
		panic("IpFinderMock.FindHostPublicIpFunc: method is nil but IpFinder.FindHostPublicIp was just called")
	}
	callInfo := struct {
	}{}
	mock.lockFindHostPublicIp.Lock()
	mock.calls.FindHostPublicIp = append(mock.calls.FindHostPublicIp, callInfo)
	mock.lockFindHostPublicIp.Unlock()
	return mock.FindHostPublicIpFunc()
}

// FindHostPublicIpCalls gets all the calls that were made to FindHostPublicIp.
// Check the length with:
//
//	len(mockedIpFinder.FindHostPublicIpCalls())
func (mock *IpFinderMock) FindHostPublicIpCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockFindHostPublicIp.RLock()
	calls = mock.calls.FindHostPublicIp
	mock.lockFindHostPublicIp.RUnlock()
	return calls
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
