// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package presentation

import (
	"sync"
)

// Ensure, that PulserMock does implement Pulser.
// If this is not the case, regenerate this file with moq.
var _ Pulser = &PulserMock{}

// PulserMock is a mock implementation of Pulser.
//
//	func TestSomethingThatUsesPulser(t *testing.T) {
//
//		// make and configure a mocked Pulser
//		mockedPulser := &PulserMock{
//			PulseFunc: func()  {
//				panic("mock out the Pulse method")
//			},
//			StartFunc: func()  {
//				panic("mock out the Start method")
//			},
//			StopFunc: func()  {
//				panic("mock out the Stop method")
//			},
//		}
//
//		// use mockedPulser in code that requires Pulser
//		// and then make assertions.
//
//	}
type PulserMock struct {
	// PulseFunc mocks the Pulse method.
	PulseFunc func()

	// StartFunc mocks the Start method.
	StartFunc func()

	// StopFunc mocks the Stop method.
	StopFunc func()

	// calls tracks calls to the methods.
	calls struct {
		// Pulse holds details about calls to the Pulse method.
		Pulse []struct {
		}
		// Start holds details about calls to the Start method.
		Start []struct {
		}
		// Stop holds details about calls to the Stop method.
		Stop []struct {
		}
	}
	lockPulse sync.RWMutex
	lockStart sync.RWMutex
	lockStop  sync.RWMutex
}

// Pulse calls PulseFunc.
func (mock *PulserMock) Pulse() {
	if mock.PulseFunc == nil {
		panic("PulserMock.PulseFunc: method is nil but Pulser.Pulse was just called")
	}
	callInfo := struct {
	}{}
	mock.lockPulse.Lock()
	mock.calls.Pulse = append(mock.calls.Pulse, callInfo)
	mock.lockPulse.Unlock()
	mock.PulseFunc()
}

// PulseCalls gets all the calls that were made to Pulse.
// Check the length with:
//
//	len(mockedPulser.PulseCalls())
func (mock *PulserMock) PulseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockPulse.RLock()
	calls = mock.calls.Pulse
	mock.lockPulse.RUnlock()
	return calls
}

// Start calls StartFunc.
func (mock *PulserMock) Start() {
	if mock.StartFunc == nil {
		panic("PulserMock.StartFunc: method is nil but Pulser.Start was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStart.Lock()
	mock.calls.Start = append(mock.calls.Start, callInfo)
	mock.lockStart.Unlock()
	mock.StartFunc()
}

// StartCalls gets all the calls that were made to Start.
// Check the length with:
//
//	len(mockedPulser.StartCalls())
func (mock *PulserMock) StartCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStart.RLock()
	calls = mock.calls.Start
	mock.lockStart.RUnlock()
	return calls
}

// Stop calls StopFunc.
func (mock *PulserMock) Stop() {
	if mock.StopFunc == nil {
		panic("PulserMock.StopFunc: method is nil but Pulser.Stop was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStop.Lock()
	mock.calls.Stop = append(mock.calls.Stop, callInfo)
	mock.lockStop.Unlock()
	mock.StopFunc()
}

// StopCalls gets all the calls that were made to Stop.
// Check the length with:
//
//	len(mockedPulser.StopCalls())
func (mock *PulserMock) StopCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStop.RLock()
	calls = mock.calls.Stop
	mock.lockStop.RUnlock()
	return calls
}