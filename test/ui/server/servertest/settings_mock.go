// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package servertest

import (
	"github.com/my-cloud/ruthenium/src/ui/server"
	"sync"
)

// Ensure, that SettingsMock does implement Settings.
// If this is not the case, regenerate this file with moq.
var _ server.Settings = &SettingsMock{}

// SettingsMock is a mock implementation of Settings.
//
//	func TestSomethingThatUsesSettings(t *testing.T) {
//
//		// make and configure a mocked Settings
//		mockedSettings := &SettingsMock{
//			HalfLifeInNanosecondsFunc: func() float64 {
//				panic("mock out the HalfLifeInNanoseconds method")
//			},
//			IncomeBaseFunc: func() uint64 {
//				panic("mock out the IncomeBase method")
//			},
//			IncomeLimitFunc: func() uint64 {
//				panic("mock out the IncomeLimit method")
//			},
//			MinimalTransactionFeeFunc: func() uint64 {
//				panic("mock out the MinimalTransactionFee method")
//			},
//			SmallestUnitsPerCoinFunc: func() uint64 {
//				panic("mock out the SmallestUnitsPerCoin method")
//			},
//			ValidationTimestampFunc: func() int64 {
//				panic("mock out the ValidationTimestamp method")
//			},
//		}
//
//		// use mockedSettings in code that requires Settings
//		// and then make assertions.
//
//	}
type SettingsMock struct {
	// HalfLifeInNanosecondsFunc mocks the HalfLifeInNanoseconds method.
	HalfLifeInNanosecondsFunc func() float64

	// IncomeBaseFunc mocks the IncomeBase method.
	IncomeBaseFunc func() uint64

	// IncomeLimitFunc mocks the IncomeLimit method.
	IncomeLimitFunc func() uint64

	// MinimalTransactionFeeFunc mocks the MinimalTransactionFee method.
	MinimalTransactionFeeFunc func() uint64

	// SmallestUnitsPerCoinFunc mocks the SmallestUnitsPerCoin method.
	SmallestUnitsPerCoinFunc func() uint64

	// ValidationTimestampFunc mocks the ValidationTimestamp method.
	ValidationTimestampFunc func() int64

	// calls tracks calls to the methods.
	calls struct {
		// HalfLifeInNanoseconds holds details about calls to the HalfLifeInNanoseconds method.
		HalfLifeInNanoseconds []struct {
		}
		// IncomeBase holds details about calls to the IncomeBase method.
		IncomeBase []struct {
		}
		// IncomeLimit holds details about calls to the IncomeLimit method.
		IncomeLimit []struct {
		}
		// MinimalTransactionFee holds details about calls to the MinimalTransactionFee method.
		MinimalTransactionFee []struct {
		}
		// SmallestUnitsPerCoin holds details about calls to the SmallestUnitsPerCoin method.
		SmallestUnitsPerCoin []struct {
		}
		// ValidationTimestamp holds details about calls to the ValidationTimestamp method.
		ValidationTimestamp []struct {
		}
	}
	lockHalfLifeInNanoseconds sync.RWMutex
	lockIncomeBase            sync.RWMutex
	lockIncomeLimit           sync.RWMutex
	lockMinimalTransactionFee sync.RWMutex
	lockSmallestUnitsPerCoin  sync.RWMutex
	lockValidationTimestamp   sync.RWMutex
}

// HalfLifeInNanoseconds calls HalfLifeInNanosecondsFunc.
func (mock *SettingsMock) HalfLifeInNanoseconds() float64 {
	if mock.HalfLifeInNanosecondsFunc == nil {
		panic("SettingsMock.HalfLifeInNanosecondsFunc: method is nil but Settings.HalfLifeInNanoseconds was just called")
	}
	callInfo := struct {
	}{}
	mock.lockHalfLifeInNanoseconds.Lock()
	mock.calls.HalfLifeInNanoseconds = append(mock.calls.HalfLifeInNanoseconds, callInfo)
	mock.lockHalfLifeInNanoseconds.Unlock()
	return mock.HalfLifeInNanosecondsFunc()
}

// HalfLifeInNanosecondsCalls gets all the calls that were made to HalfLifeInNanoseconds.
// Check the length with:
//
//	len(mockedSettings.HalfLifeInNanosecondsCalls())
func (mock *SettingsMock) HalfLifeInNanosecondsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockHalfLifeInNanoseconds.RLock()
	calls = mock.calls.HalfLifeInNanoseconds
	mock.lockHalfLifeInNanoseconds.RUnlock()
	return calls
}

// IncomeBase calls IncomeBaseFunc.
func (mock *SettingsMock) IncomeBase() uint64 {
	if mock.IncomeBaseFunc == nil {
		panic("SettingsMock.IncomeBaseFunc: method is nil but Settings.IncomeBase was just called")
	}
	callInfo := struct {
	}{}
	mock.lockIncomeBase.Lock()
	mock.calls.IncomeBase = append(mock.calls.IncomeBase, callInfo)
	mock.lockIncomeBase.Unlock()
	return mock.IncomeBaseFunc()
}

// IncomeBaseCalls gets all the calls that were made to IncomeBase.
// Check the length with:
//
//	len(mockedSettings.IncomeBaseCalls())
func (mock *SettingsMock) IncomeBaseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeBase.RLock()
	calls = mock.calls.IncomeBase
	mock.lockIncomeBase.RUnlock()
	return calls
}

// IncomeLimit calls IncomeLimitFunc.
func (mock *SettingsMock) IncomeLimit() uint64 {
	if mock.IncomeLimitFunc == nil {
		panic("SettingsMock.IncomeLimitFunc: method is nil but Settings.IncomeLimit was just called")
	}
	callInfo := struct {
	}{}
	mock.lockIncomeLimit.Lock()
	mock.calls.IncomeLimit = append(mock.calls.IncomeLimit, callInfo)
	mock.lockIncomeLimit.Unlock()
	return mock.IncomeLimitFunc()
}

// IncomeLimitCalls gets all the calls that were made to IncomeLimit.
// Check the length with:
//
//	len(mockedSettings.IncomeLimitCalls())
func (mock *SettingsMock) IncomeLimitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeLimit.RLock()
	calls = mock.calls.IncomeLimit
	mock.lockIncomeLimit.RUnlock()
	return calls
}

// MinimalTransactionFee calls MinimalTransactionFeeFunc.
func (mock *SettingsMock) MinimalTransactionFee() uint64 {
	if mock.MinimalTransactionFeeFunc == nil {
		panic("SettingsMock.MinimalTransactionFeeFunc: method is nil but Settings.MinimalTransactionFee was just called")
	}
	callInfo := struct {
	}{}
	mock.lockMinimalTransactionFee.Lock()
	mock.calls.MinimalTransactionFee = append(mock.calls.MinimalTransactionFee, callInfo)
	mock.lockMinimalTransactionFee.Unlock()
	return mock.MinimalTransactionFeeFunc()
}

// MinimalTransactionFeeCalls gets all the calls that were made to MinimalTransactionFee.
// Check the length with:
//
//	len(mockedSettings.MinimalTransactionFeeCalls())
func (mock *SettingsMock) MinimalTransactionFeeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockMinimalTransactionFee.RLock()
	calls = mock.calls.MinimalTransactionFee
	mock.lockMinimalTransactionFee.RUnlock()
	return calls
}

// SmallestUnitsPerCoin calls SmallestUnitsPerCoinFunc.
func (mock *SettingsMock) SmallestUnitsPerCoin() uint64 {
	if mock.SmallestUnitsPerCoinFunc == nil {
		panic("SettingsMock.SmallestUnitsPerCoinFunc: method is nil but Settings.SmallestUnitsPerCoin was just called")
	}
	callInfo := struct {
	}{}
	mock.lockSmallestUnitsPerCoin.Lock()
	mock.calls.SmallestUnitsPerCoin = append(mock.calls.SmallestUnitsPerCoin, callInfo)
	mock.lockSmallestUnitsPerCoin.Unlock()
	return mock.SmallestUnitsPerCoinFunc()
}

// SmallestUnitsPerCoinCalls gets all the calls that were made to SmallestUnitsPerCoin.
// Check the length with:
//
//	len(mockedSettings.SmallestUnitsPerCoinCalls())
func (mock *SettingsMock) SmallestUnitsPerCoinCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSmallestUnitsPerCoin.RLock()
	calls = mock.calls.SmallestUnitsPerCoin
	mock.lockSmallestUnitsPerCoin.RUnlock()
	return calls
}

// ValidationTimestamp calls ValidationTimestampFunc.
func (mock *SettingsMock) ValidationTimestamp() int64 {
	if mock.ValidationTimestampFunc == nil {
		panic("SettingsMock.ValidationTimestampFunc: method is nil but Settings.ValidationTimestamp was just called")
	}
	callInfo := struct {
	}{}
	mock.lockValidationTimestamp.Lock()
	mock.calls.ValidationTimestamp = append(mock.calls.ValidationTimestamp, callInfo)
	mock.lockValidationTimestamp.Unlock()
	return mock.ValidationTimestampFunc()
}

// ValidationTimestampCalls gets all the calls that were made to ValidationTimestamp.
// Check the length with:
//
//	len(mockedSettings.ValidationTimestampCalls())
func (mock *SettingsMock) ValidationTimestampCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidationTimestamp.RLock()
	calls = mock.calls.ValidationTimestamp
	mock.lockValidationTimestamp.RUnlock()
	return calls
}
