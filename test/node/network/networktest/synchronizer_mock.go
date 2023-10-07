// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package networktest

import (
	"github.com/my-cloud/ruthenium/src/node/network"
	"sync"
)

// Ensure, that SynchronizerMock does implement Synchronizer.
// If this is not the case, regenerate this file with moq.
var _ network.Synchronizer = &SynchronizerMock{}

// SynchronizerMock is a mock implementation of Synchronizer.
//
//	func TestSomethingThatUsesSynchronizer(t *testing.T) {
//
//		// make and configure a mocked Synchronizer
//		mockedSynchronizer := &SynchronizerMock{
//			AddTargetsFunc: func(targets []string)  {
//				panic("mock out the AddTargets method")
//			},
//			HostTargetFunc: func() string {
//				panic("mock out the HostTarget method")
//			},
//			IncentiveFunc: func(target string)  {
//				panic("mock out the Incentive method")
//			},
//			NeighborsFunc: func() []Neighbor {
//				panic("mock out the Neighbors method")
//			},
//		}
//
//		// use mockedSynchronizer in code that requires Synchronizer
//		// and then make assertions.
//
//	}
type SynchronizerMock struct {
	// AddTargetsFunc mocks the AddTargets method.
	AddTargetsFunc func(targets []string)

	// HostTargetFunc mocks the HostTarget method.
	HostTargetFunc func() string

	// IncentiveFunc mocks the Incentive method.
	IncentiveFunc func(target string)

	// NeighborsFunc mocks the Neighbors method.
	NeighborsFunc func() []network.Neighbor

	// calls tracks calls to the methods.
	calls struct {
		// AddTargets holds details about calls to the AddTargets method.
		AddTargets []struct {
			// Targets is the targets argument value.
			Targets []string
		}
		// HostTarget holds details about calls to the HostTarget method.
		HostTarget []struct {
		}
		// Incentive holds details about calls to the Incentive method.
		Incentive []struct {
			// Target is the target argument value.
			Target string
		}
		// Neighbors holds details about calls to the Neighbors method.
		Neighbors []struct {
		}
	}
	lockAddTargets sync.RWMutex
	lockHostTarget sync.RWMutex
	lockIncentive  sync.RWMutex
	lockNeighbors  sync.RWMutex
}

// AddTargets calls AddTargetsFunc.
func (mock *SynchronizerMock) AddTargets(targets []string) {
	if mock.AddTargetsFunc == nil {
		panic("SynchronizerMock.AddTargetsFunc: method is nil but Synchronizer.AddTargets was just called")
	}
	callInfo := struct {
		Targets []string
	}{
		Targets: targets,
	}
	mock.lockAddTargets.Lock()
	mock.calls.AddTargets = append(mock.calls.AddTargets, callInfo)
	mock.lockAddTargets.Unlock()
	mock.AddTargetsFunc(targets)
}

// AddTargetsCalls gets all the calls that were made to AddTargets.
// Check the length with:
//
//	len(mockedSynchronizer.AddTargetsCalls())
func (mock *SynchronizerMock) AddTargetsCalls() []struct {
	Targets []string
} {
	var calls []struct {
		Targets []string
	}
	mock.lockAddTargets.RLock()
	calls = mock.calls.AddTargets
	mock.lockAddTargets.RUnlock()
	return calls
}

// HostTarget calls HostTargetFunc.
func (mock *SynchronizerMock) HostTarget() string {
	if mock.HostTargetFunc == nil {
		panic("SynchronizerMock.HostTargetFunc: method is nil but Synchronizer.HostTarget was just called")
	}
	callInfo := struct {
	}{}
	mock.lockHostTarget.Lock()
	mock.calls.HostTarget = append(mock.calls.HostTarget, callInfo)
	mock.lockHostTarget.Unlock()
	return mock.HostTargetFunc()
}

// HostTargetCalls gets all the calls that were made to HostTarget.
// Check the length with:
//
//	len(mockedSynchronizer.HostTargetCalls())
func (mock *SynchronizerMock) HostTargetCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockHostTarget.RLock()
	calls = mock.calls.HostTarget
	mock.lockHostTarget.RUnlock()
	return calls
}

// Incentive calls IncentiveFunc.
func (mock *SynchronizerMock) Incentive(target string) {
	if mock.IncentiveFunc == nil {
		panic("SynchronizerMock.IncentiveFunc: method is nil but Synchronizer.Incentive was just called")
	}
	callInfo := struct {
		Target string
	}{
		Target: target,
	}
	mock.lockIncentive.Lock()
	mock.calls.Incentive = append(mock.calls.Incentive, callInfo)
	mock.lockIncentive.Unlock()
	mock.IncentiveFunc(target)
}

// IncentiveCalls gets all the calls that were made to Incentive.
// Check the length with:
//
//	len(mockedSynchronizer.IncentiveCalls())
func (mock *SynchronizerMock) IncentiveCalls() []struct {
	Target string
} {
	var calls []struct {
		Target string
	}
	mock.lockIncentive.RLock()
	calls = mock.calls.Incentive
	mock.lockIncentive.RUnlock()
	return calls
}

// Neighbors calls NeighborsFunc.
func (mock *SynchronizerMock) Neighbors() []network.Neighbor {
	if mock.NeighborsFunc == nil {
		panic("SynchronizerMock.NeighborsFunc: method is nil but Synchronizer.Neighbors was just called")
	}
	callInfo := struct {
	}{}
	mock.lockNeighbors.Lock()
	mock.calls.Neighbors = append(mock.calls.Neighbors, callInfo)
	mock.lockNeighbors.Unlock()
	return mock.NeighborsFunc()
}

// NeighborsCalls gets all the calls that were made to Neighbors.
// Check the length with:
//
//	len(mockedSynchronizer.NeighborsCalls())
func (mock *SynchronizerMock) NeighborsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockNeighbors.RLock()
	calls = mock.calls.Neighbors
	mock.lockNeighbors.RUnlock()
	return calls
}
