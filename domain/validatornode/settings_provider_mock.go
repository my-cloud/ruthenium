// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package validatornode

import (
	"sync"
	"time"
)

// Ensure, that SettingsProviderMock does implement SettingsProvider.
// If this is not the case, regenerate this file with moq.
var _ SettingsProvider = &SettingsProviderMock{}

// SettingsProviderMock is a mock implementation of SettingsProvider.
//
//	func TestSomethingThatUsesSettingsProvider(t *testing.T) {
//
//		// make and configure a mocked SettingsProvider
//		mockedSettingsProvider := &SettingsProviderMock{
//			BlocksCountLimitFunc: func() uint64 {
//				panic("mock out the BlocksCountLimit method")
//			},
//			GenesisAmountFunc: func() uint64 {
//				panic("mock out the GenesisAmount method")
//			},
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
//			ValidationTimeoutFunc: func() time.Duration {
//				panic("mock out the ValidationTimeout method")
//			},
//			ValidationTimestampFunc: func() int64 {
//				panic("mock out the ValidationTimestamp method")
//			},
//		}
//
//		// use mockedSettingsProvider in code that requires SettingsProvider
//		// and then make assertions.
//
//	}
type SettingsProviderMock struct {
	// BlocksCountLimitFunc mocks the BlocksCountLimit method.
	BlocksCountLimitFunc func() uint64

	// GenesisAmountFunc mocks the GenesisAmount method.
	GenesisAmountFunc func() uint64

	// HalfLifeInNanosecondsFunc mocks the HalfLifeInNanoseconds method.
	HalfLifeInNanosecondsFunc func() float64

	// IncomeBaseFunc mocks the IncomeBase method.
	IncomeBaseFunc func() uint64

	// IncomeLimitFunc mocks the IncomeLimit method.
	IncomeLimitFunc func() uint64

	// MinimalTransactionFeeFunc mocks the MinimalTransactionFee method.
	MinimalTransactionFeeFunc func() uint64

	// ValidationTimeoutFunc mocks the ValidationTimeout method.
	ValidationTimeoutFunc func() time.Duration

	// ValidationTimestampFunc mocks the ValidationTimestamp method.
	ValidationTimestampFunc func() int64

	// calls tracks calls to the methods.
	calls struct {
		// BlocksCountLimit holds details about calls to the BlocksCountLimit method.
		BlocksCountLimit []struct {
		}
		// GenesisAmount holds details about calls to the GenesisAmount method.
		GenesisAmount []struct {
		}
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
		// ValidationTimeout holds details about calls to the ValidationTimeout method.
		ValidationTimeout []struct {
		}
		// ValidationTimestamp holds details about calls to the ValidationTimestamp method.
		ValidationTimestamp []struct {
		}
	}
	lockBlocksCountLimit      sync.RWMutex
	lockGenesisAmount         sync.RWMutex
	lockHalfLifeInNanoseconds sync.RWMutex
	lockIncomeBase            sync.RWMutex
	lockIncomeLimit           sync.RWMutex
	lockMinimalTransactionFee sync.RWMutex
	lockValidationTimeout     sync.RWMutex
	lockValidationTimestamp   sync.RWMutex
}

// BlocksCountLimit calls BlocksCountLimitFunc.
func (mock *SettingsProviderMock) BlocksCountLimit() uint64 {
	if mock.BlocksCountLimitFunc == nil {
		panic("SettingsProviderMock.BlocksCountLimitFunc: method is nil but SettingsProvider.BlocksCountLimit was just called")
	}
	callInfo := struct {
	}{}
	mock.lockBlocksCountLimit.Lock()
	mock.calls.BlocksCountLimit = append(mock.calls.BlocksCountLimit, callInfo)
	mock.lockBlocksCountLimit.Unlock()
	return mock.BlocksCountLimitFunc()
}

// BlocksCountLimitCalls gets all the calls that were made to BlocksCountLimit.
// Check the length with:
//
//	len(mockedSettingsProvider.BlocksCountLimitCalls())
func (mock *SettingsProviderMock) BlocksCountLimitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockBlocksCountLimit.RLock()
	calls = mock.calls.BlocksCountLimit
	mock.lockBlocksCountLimit.RUnlock()
	return calls
}

// GenesisAmount calls GenesisAmountFunc.
func (mock *SettingsProviderMock) GenesisAmount() uint64 {
	if mock.GenesisAmountFunc == nil {
		panic("SettingsProviderMock.GenesisAmountFunc: method is nil but SettingsProvider.GenesisAmount was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGenesisAmount.Lock()
	mock.calls.GenesisAmount = append(mock.calls.GenesisAmount, callInfo)
	mock.lockGenesisAmount.Unlock()
	return mock.GenesisAmountFunc()
}

// GenesisAmountCalls gets all the calls that were made to GenesisAmount.
// Check the length with:
//
//	len(mockedSettingsProvider.GenesisAmountCalls())
func (mock *SettingsProviderMock) GenesisAmountCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGenesisAmount.RLock()
	calls = mock.calls.GenesisAmount
	mock.lockGenesisAmount.RUnlock()
	return calls
}

// HalfLifeInNanoseconds calls HalfLifeInNanosecondsFunc.
func (mock *SettingsProviderMock) HalfLifeInNanoseconds() float64 {
	if mock.HalfLifeInNanosecondsFunc == nil {
		panic("SettingsProviderMock.HalfLifeInNanosecondsFunc: method is nil but SettingsProvider.HalfLifeInNanoseconds was just called")
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
//	len(mockedSettingsProvider.HalfLifeInNanosecondsCalls())
func (mock *SettingsProviderMock) HalfLifeInNanosecondsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockHalfLifeInNanoseconds.RLock()
	calls = mock.calls.HalfLifeInNanoseconds
	mock.lockHalfLifeInNanoseconds.RUnlock()
	return calls
}

// IncomeBase calls IncomeBaseFunc.
func (mock *SettingsProviderMock) IncomeBase() uint64 {
	if mock.IncomeBaseFunc == nil {
		panic("SettingsProviderMock.IncomeBaseFunc: method is nil but SettingsProvider.IncomeBase was just called")
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
//	len(mockedSettingsProvider.IncomeBaseCalls())
func (mock *SettingsProviderMock) IncomeBaseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeBase.RLock()
	calls = mock.calls.IncomeBase
	mock.lockIncomeBase.RUnlock()
	return calls
}

// IncomeLimit calls IncomeLimitFunc.
func (mock *SettingsProviderMock) IncomeLimit() uint64 {
	if mock.IncomeLimitFunc == nil {
		panic("SettingsProviderMock.IncomeLimitFunc: method is nil but SettingsProvider.IncomeLimit was just called")
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
//	len(mockedSettingsProvider.IncomeLimitCalls())
func (mock *SettingsProviderMock) IncomeLimitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeLimit.RLock()
	calls = mock.calls.IncomeLimit
	mock.lockIncomeLimit.RUnlock()
	return calls
}

// MinimalTransactionFee calls MinimalTransactionFeeFunc.
func (mock *SettingsProviderMock) MinimalTransactionFee() uint64 {
	if mock.MinimalTransactionFeeFunc == nil {
		panic("SettingsProviderMock.MinimalTransactionFeeFunc: method is nil but SettingsProvider.MinimalTransactionFee was just called")
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
//	len(mockedSettingsProvider.MinimalTransactionFeeCalls())
func (mock *SettingsProviderMock) MinimalTransactionFeeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockMinimalTransactionFee.RLock()
	calls = mock.calls.MinimalTransactionFee
	mock.lockMinimalTransactionFee.RUnlock()
	return calls
}

// ValidationTimeout calls ValidationTimeoutFunc.
func (mock *SettingsProviderMock) ValidationTimeout() time.Duration {
	if mock.ValidationTimeoutFunc == nil {
		panic("SettingsProviderMock.ValidationTimeoutFunc: method is nil but SettingsProvider.ValidationTimeout was just called")
	}
	callInfo := struct {
	}{}
	mock.lockValidationTimeout.Lock()
	mock.calls.ValidationTimeout = append(mock.calls.ValidationTimeout, callInfo)
	mock.lockValidationTimeout.Unlock()
	return mock.ValidationTimeoutFunc()
}

// ValidationTimeoutCalls gets all the calls that were made to ValidationTimeout.
// Check the length with:
//
//	len(mockedSettingsProvider.ValidationTimeoutCalls())
func (mock *SettingsProviderMock) ValidationTimeoutCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidationTimeout.RLock()
	calls = mock.calls.ValidationTimeout
	mock.lockValidationTimeout.RUnlock()
	return calls
}

// ValidationTimestamp calls ValidationTimestampFunc.
func (mock *SettingsProviderMock) ValidationTimestamp() int64 {
	if mock.ValidationTimestampFunc == nil {
		panic("SettingsProviderMock.ValidationTimestampFunc: method is nil but SettingsProvider.ValidationTimestamp was just called")
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
//	len(mockedSettingsProvider.ValidationTimestampCalls())
func (mock *SettingsProviderMock) ValidationTimestampCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidationTimestamp.RLock()
	calls = mock.calls.ValidationTimestamp
	mock.lockValidationTimestamp.RUnlock()
	return calls
}