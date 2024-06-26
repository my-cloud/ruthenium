// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package application

import (
	"sync"
	"time"
)

// Ensure, that ProtocolSettingsProviderMock does implement ProtocolSettingsProvider.
// If this is not the case, regenerate this file with moq.
var _ ProtocolSettingsProvider = &ProtocolSettingsProviderMock{}

// ProtocolSettingsProviderMock is a mock implementation of ProtocolSettingsProvider.
//
//	func TestSomethingThatUsesProtocolSettingsProvider(t *testing.T) {
//
//		// make and configure a mocked ProtocolSettingsProvider
//		mockedProtocolSettingsProvider := &ProtocolSettingsProviderMock{
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
//			SmallestUnitsPerCoinFunc: func() uint64 {
//				panic("mock out the SmallestUnitsPerCoin method")
//			},
//			ValidationTimeoutFunc: func() time.Duration {
//				panic("mock out the ValidationTimeout method")
//			},
//			ValidationTimerFunc: func() time.Duration {
//				panic("mock out the ValidationTimer method")
//			},
//			ValidationTimestampFunc: func() int64 {
//				panic("mock out the ValidationTimestamp method")
//			},
//			VerificationsCountPerValidationFunc: func() int64 {
//				panic("mock out the VerificationsCountPerValidation method")
//			},
//		}
//
//		// use mockedProtocolSettingsProvider in code that requires ProtocolSettingsProvider
//		// and then make assertions.
//
//	}
type ProtocolSettingsProviderMock struct {
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

	// SmallestUnitsPerCoinFunc mocks the SmallestUnitsPerCoin method.
	SmallestUnitsPerCoinFunc func() uint64

	// ValidationTimeoutFunc mocks the ValidationTimeout method.
	ValidationTimeoutFunc func() time.Duration

	// ValidationTimerFunc mocks the ValidationTimer method.
	ValidationTimerFunc func() time.Duration

	// ValidationTimestampFunc mocks the ValidationTimestamp method.
	ValidationTimestampFunc func() int64

	// VerificationsCountPerValidationFunc mocks the VerificationsCountPerValidation method.
	VerificationsCountPerValidationFunc func() int64

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
		// SmallestUnitsPerCoin holds details about calls to the SmallestUnitsPerCoin method.
		SmallestUnitsPerCoin []struct {
		}
		// ValidationTimeout holds details about calls to the ValidationTimeout method.
		ValidationTimeout []struct {
		}
		// ValidationTimer holds details about calls to the ValidationTimer method.
		ValidationTimer []struct {
		}
		// ValidationTimestamp holds details about calls to the ValidationTimestamp method.
		ValidationTimestamp []struct {
		}
		// VerificationsCountPerValidation holds details about calls to the VerificationsCountPerValidation method.
		VerificationsCountPerValidation []struct {
		}
	}
	lockBlocksCountLimit                sync.RWMutex
	lockGenesisAmount                   sync.RWMutex
	lockHalfLifeInNanoseconds           sync.RWMutex
	lockIncomeBase                      sync.RWMutex
	lockIncomeLimit                     sync.RWMutex
	lockMinimalTransactionFee           sync.RWMutex
	lockSmallestUnitsPerCoin            sync.RWMutex
	lockValidationTimeout               sync.RWMutex
	lockValidationTimer                 sync.RWMutex
	lockValidationTimestamp             sync.RWMutex
	lockVerificationsCountPerValidation sync.RWMutex
}

// BlocksCountLimit calls BlocksCountLimitFunc.
func (mock *ProtocolSettingsProviderMock) BlocksCountLimit() uint64 {
	if mock.BlocksCountLimitFunc == nil {
		panic("ProtocolSettingsProviderMock.BlocksCountLimitFunc: method is nil but ProtocolSettingsProvider.BlocksCountLimit was just called")
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
//	len(mockedProtocolSettingsProvider.BlocksCountLimitCalls())
func (mock *ProtocolSettingsProviderMock) BlocksCountLimitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockBlocksCountLimit.RLock()
	calls = mock.calls.BlocksCountLimit
	mock.lockBlocksCountLimit.RUnlock()
	return calls
}

// GenesisAmount calls GenesisAmountFunc.
func (mock *ProtocolSettingsProviderMock) GenesisAmount() uint64 {
	if mock.GenesisAmountFunc == nil {
		panic("ProtocolSettingsProviderMock.GenesisAmountFunc: method is nil but ProtocolSettingsProvider.GenesisAmount was just called")
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
//	len(mockedProtocolSettingsProvider.GenesisAmountCalls())
func (mock *ProtocolSettingsProviderMock) GenesisAmountCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGenesisAmount.RLock()
	calls = mock.calls.GenesisAmount
	mock.lockGenesisAmount.RUnlock()
	return calls
}

// HalfLifeInNanoseconds calls HalfLifeInNanosecondsFunc.
func (mock *ProtocolSettingsProviderMock) HalfLifeInNanoseconds() float64 {
	if mock.HalfLifeInNanosecondsFunc == nil {
		panic("ProtocolSettingsProviderMock.HalfLifeInNanosecondsFunc: method is nil but ProtocolSettingsProvider.HalfLifeInNanoseconds was just called")
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
//	len(mockedProtocolSettingsProvider.HalfLifeInNanosecondsCalls())
func (mock *ProtocolSettingsProviderMock) HalfLifeInNanosecondsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockHalfLifeInNanoseconds.RLock()
	calls = mock.calls.HalfLifeInNanoseconds
	mock.lockHalfLifeInNanoseconds.RUnlock()
	return calls
}

// IncomeBase calls IncomeBaseFunc.
func (mock *ProtocolSettingsProviderMock) IncomeBase() uint64 {
	if mock.IncomeBaseFunc == nil {
		panic("ProtocolSettingsProviderMock.IncomeBaseFunc: method is nil but ProtocolSettingsProvider.IncomeBase was just called")
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
//	len(mockedProtocolSettingsProvider.IncomeBaseCalls())
func (mock *ProtocolSettingsProviderMock) IncomeBaseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeBase.RLock()
	calls = mock.calls.IncomeBase
	mock.lockIncomeBase.RUnlock()
	return calls
}

// IncomeLimit calls IncomeLimitFunc.
func (mock *ProtocolSettingsProviderMock) IncomeLimit() uint64 {
	if mock.IncomeLimitFunc == nil {
		panic("ProtocolSettingsProviderMock.IncomeLimitFunc: method is nil but ProtocolSettingsProvider.IncomeLimit was just called")
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
//	len(mockedProtocolSettingsProvider.IncomeLimitCalls())
func (mock *ProtocolSettingsProviderMock) IncomeLimitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeLimit.RLock()
	calls = mock.calls.IncomeLimit
	mock.lockIncomeLimit.RUnlock()
	return calls
}

// MinimalTransactionFee calls MinimalTransactionFeeFunc.
func (mock *ProtocolSettingsProviderMock) MinimalTransactionFee() uint64 {
	if mock.MinimalTransactionFeeFunc == nil {
		panic("ProtocolSettingsProviderMock.MinimalTransactionFeeFunc: method is nil but ProtocolSettingsProvider.MinimalTransactionFee was just called")
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
//	len(mockedProtocolSettingsProvider.MinimalTransactionFeeCalls())
func (mock *ProtocolSettingsProviderMock) MinimalTransactionFeeCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockMinimalTransactionFee.RLock()
	calls = mock.calls.MinimalTransactionFee
	mock.lockMinimalTransactionFee.RUnlock()
	return calls
}

// SmallestUnitsPerCoin calls SmallestUnitsPerCoinFunc.
func (mock *ProtocolSettingsProviderMock) SmallestUnitsPerCoin() uint64 {
	if mock.SmallestUnitsPerCoinFunc == nil {
		panic("ProtocolSettingsProviderMock.SmallestUnitsPerCoinFunc: method is nil but ProtocolSettingsProvider.SmallestUnitsPerCoin was just called")
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
//	len(mockedProtocolSettingsProvider.SmallestUnitsPerCoinCalls())
func (mock *ProtocolSettingsProviderMock) SmallestUnitsPerCoinCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockSmallestUnitsPerCoin.RLock()
	calls = mock.calls.SmallestUnitsPerCoin
	mock.lockSmallestUnitsPerCoin.RUnlock()
	return calls
}

// ValidationTimeout calls ValidationTimeoutFunc.
func (mock *ProtocolSettingsProviderMock) ValidationTimeout() time.Duration {
	if mock.ValidationTimeoutFunc == nil {
		panic("ProtocolSettingsProviderMock.ValidationTimeoutFunc: method is nil but ProtocolSettingsProvider.ValidationTimeout was just called")
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
//	len(mockedProtocolSettingsProvider.ValidationTimeoutCalls())
func (mock *ProtocolSettingsProviderMock) ValidationTimeoutCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidationTimeout.RLock()
	calls = mock.calls.ValidationTimeout
	mock.lockValidationTimeout.RUnlock()
	return calls
}

// ValidationTimer calls ValidationTimerFunc.
func (mock *ProtocolSettingsProviderMock) ValidationTimer() time.Duration {
	if mock.ValidationTimerFunc == nil {
		panic("ProtocolSettingsProviderMock.ValidationTimerFunc: method is nil but ProtocolSettingsProvider.ValidationTimer was just called")
	}
	callInfo := struct {
	}{}
	mock.lockValidationTimer.Lock()
	mock.calls.ValidationTimer = append(mock.calls.ValidationTimer, callInfo)
	mock.lockValidationTimer.Unlock()
	return mock.ValidationTimerFunc()
}

// ValidationTimerCalls gets all the calls that were made to ValidationTimer.
// Check the length with:
//
//	len(mockedProtocolSettingsProvider.ValidationTimerCalls())
func (mock *ProtocolSettingsProviderMock) ValidationTimerCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidationTimer.RLock()
	calls = mock.calls.ValidationTimer
	mock.lockValidationTimer.RUnlock()
	return calls
}

// ValidationTimestamp calls ValidationTimestampFunc.
func (mock *ProtocolSettingsProviderMock) ValidationTimestamp() int64 {
	if mock.ValidationTimestampFunc == nil {
		panic("ProtocolSettingsProviderMock.ValidationTimestampFunc: method is nil but ProtocolSettingsProvider.ValidationTimestamp was just called")
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
//	len(mockedProtocolSettingsProvider.ValidationTimestampCalls())
func (mock *ProtocolSettingsProviderMock) ValidationTimestampCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidationTimestamp.RLock()
	calls = mock.calls.ValidationTimestamp
	mock.lockValidationTimestamp.RUnlock()
	return calls
}

// VerificationsCountPerValidation calls VerificationsCountPerValidationFunc.
func (mock *ProtocolSettingsProviderMock) VerificationsCountPerValidation() int64 {
	if mock.VerificationsCountPerValidationFunc == nil {
		panic("ProtocolSettingsProviderMock.VerificationsCountPerValidationFunc: method is nil but ProtocolSettingsProvider.VerificationsCountPerValidation was just called")
	}
	callInfo := struct {
	}{}
	mock.lockVerificationsCountPerValidation.Lock()
	mock.calls.VerificationsCountPerValidation = append(mock.calls.VerificationsCountPerValidation, callInfo)
	mock.lockVerificationsCountPerValidation.Unlock()
	return mock.VerificationsCountPerValidationFunc()
}

// VerificationsCountPerValidationCalls gets all the calls that were made to VerificationsCountPerValidation.
// Check the length with:
//
//	len(mockedProtocolSettingsProvider.VerificationsCountPerValidationCalls())
func (mock *ProtocolSettingsProviderMock) VerificationsCountPerValidationCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockVerificationsCountPerValidation.RLock()
	calls = mock.calls.VerificationsCountPerValidation
	mock.lockVerificationsCountPerValidation.RUnlock()
	return calls
}
