// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package protocoltest

import (
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"sync"
	"time"
)

// Ensure, that SettingsMock does implement Settings.
// If this is not the case, regenerate this file with moq.
var _ protocol.Settings = &SettingsMock{}

// SettingsMock is a mock implementation of Settings.
//
//	func TestSomethingThatUsesSettings(t *testing.T) {
//
//		// make and configure a mocked Settings
//		mockedSettings := &SettingsMock{
//			BlocksCountLimitFunc: func() uint64 {
//				panic("mock out the BlocksCountLimit method")
//			},
//			GenesisAmountInParticlesFunc: func() uint64 {
//				panic("mock out the GenesisAmountInParticles method")
//			},
//			HalfLifeInNanosecondsFunc: func() float64 {
//				panic("mock out the HalfLifeInNanoseconds method")
//			},
//			IncomeBaseInParticlesFunc: func() uint64 {
//				panic("mock out the IncomeBaseInParticles method")
//			},
//			IncomeLimitInParticlesFunc: func() uint64 {
//				panic("mock out the IncomeLimitInParticles method")
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
//		// use mockedSettings in code that requires Settings
//		// and then make assertions.
//
//	}
type SettingsMock struct {
	// BlocksCountLimitFunc mocks the BlocksCountLimit method.
	BlocksCountLimitFunc func() uint64

	// GenesisAmountInParticlesFunc mocks the GenesisAmountInParticles method.
	GenesisAmountInParticlesFunc func() uint64

	// HalfLifeInNanosecondsFunc mocks the HalfLifeInNanoseconds method.
	HalfLifeInNanosecondsFunc func() float64

	// IncomeBaseInParticlesFunc mocks the IncomeBaseInParticles method.
	IncomeBaseInParticlesFunc func() uint64

	// IncomeLimitInParticlesFunc mocks the IncomeLimitInParticles method.
	IncomeLimitInParticlesFunc func() uint64

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
		// GenesisAmountInParticles holds details about calls to the GenesisAmountInParticles method.
		GenesisAmountInParticles []struct {
		}
		// HalfLifeInNanoseconds holds details about calls to the HalfLifeInNanoseconds method.
		HalfLifeInNanoseconds []struct {
		}
		// IncomeBaseInParticles holds details about calls to the IncomeBaseInParticles method.
		IncomeBaseInParticles []struct {
		}
		// IncomeLimitInParticles holds details about calls to the IncomeLimitInParticles method.
		IncomeLimitInParticles []struct {
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
	lockBlocksCountLimit         sync.RWMutex
	lockGenesisAmountInParticles sync.RWMutex
	lockHalfLifeInNanoseconds    sync.RWMutex
	lockIncomeBaseInParticles    sync.RWMutex
	lockIncomeLimitInParticles   sync.RWMutex
	lockMinimalTransactionFee    sync.RWMutex
	lockValidationTimeout        sync.RWMutex
	lockValidationTimestamp      sync.RWMutex
}

// BlocksCountLimit calls BlocksCountLimitFunc.
func (mock *SettingsMock) BlocksCountLimit() uint64 {
	if mock.BlocksCountLimitFunc == nil {
		panic("SettingsMock.BlocksCountLimitFunc: method is nil but Settings.BlocksCountLimit was just called")
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
//	len(mockedSettings.BlocksCountLimitCalls())
func (mock *SettingsMock) BlocksCountLimitCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockBlocksCountLimit.RLock()
	calls = mock.calls.BlocksCountLimit
	mock.lockBlocksCountLimit.RUnlock()
	return calls
}

// GenesisAmountInParticles calls GenesisAmountInParticlesFunc.
func (mock *SettingsMock) GenesisAmountInParticles() uint64 {
	if mock.GenesisAmountInParticlesFunc == nil {
		panic("SettingsMock.GenesisAmountInParticlesFunc: method is nil but Settings.GenesisAmountInParticles was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGenesisAmountInParticles.Lock()
	mock.calls.GenesisAmountInParticles = append(mock.calls.GenesisAmountInParticles, callInfo)
	mock.lockGenesisAmountInParticles.Unlock()
	return mock.GenesisAmountInParticlesFunc()
}

// GenesisAmountInParticlesCalls gets all the calls that were made to GenesisAmountInParticles.
// Check the length with:
//
//	len(mockedSettings.GenesisAmountInParticlesCalls())
func (mock *SettingsMock) GenesisAmountInParticlesCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGenesisAmountInParticles.RLock()
	calls = mock.calls.GenesisAmountInParticles
	mock.lockGenesisAmountInParticles.RUnlock()
	return calls
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

// IncomeBaseInParticles calls IncomeBaseInParticlesFunc.
func (mock *SettingsMock) IncomeBaseInParticles() uint64 {
	if mock.IncomeBaseInParticlesFunc == nil {
		panic("SettingsMock.IncomeBaseInParticlesFunc: method is nil but Settings.IncomeBaseInParticles was just called")
	}
	callInfo := struct {
	}{}
	mock.lockIncomeBaseInParticles.Lock()
	mock.calls.IncomeBaseInParticles = append(mock.calls.IncomeBaseInParticles, callInfo)
	mock.lockIncomeBaseInParticles.Unlock()
	return mock.IncomeBaseInParticlesFunc()
}

// IncomeBaseInParticlesCalls gets all the calls that were made to IncomeBaseInParticles.
// Check the length with:
//
//	len(mockedSettings.IncomeBaseInParticlesCalls())
func (mock *SettingsMock) IncomeBaseInParticlesCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeBaseInParticles.RLock()
	calls = mock.calls.IncomeBaseInParticles
	mock.lockIncomeBaseInParticles.RUnlock()
	return calls
}

// IncomeLimitInParticles calls IncomeLimitInParticlesFunc.
func (mock *SettingsMock) IncomeLimitInParticles() uint64 {
	if mock.IncomeLimitInParticlesFunc == nil {
		panic("SettingsMock.IncomeLimitInParticlesFunc: method is nil but Settings.IncomeLimitInParticles was just called")
	}
	callInfo := struct {
	}{}
	mock.lockIncomeLimitInParticles.Lock()
	mock.calls.IncomeLimitInParticles = append(mock.calls.IncomeLimitInParticles, callInfo)
	mock.lockIncomeLimitInParticles.Unlock()
	return mock.IncomeLimitInParticlesFunc()
}

// IncomeLimitInParticlesCalls gets all the calls that were made to IncomeLimitInParticles.
// Check the length with:
//
//	len(mockedSettings.IncomeLimitInParticlesCalls())
func (mock *SettingsMock) IncomeLimitInParticlesCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIncomeLimitInParticles.RLock()
	calls = mock.calls.IncomeLimitInParticles
	mock.lockIncomeLimitInParticles.RUnlock()
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

// ValidationTimeout calls ValidationTimeoutFunc.
func (mock *SettingsMock) ValidationTimeout() time.Duration {
	if mock.ValidationTimeoutFunc == nil {
		panic("SettingsMock.ValidationTimeoutFunc: method is nil but Settings.ValidationTimeout was just called")
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
//	len(mockedSettings.ValidationTimeoutCalls())
func (mock *SettingsMock) ValidationTimeoutCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidationTimeout.RLock()
	calls = mock.calls.ValidationTimeout
	mock.lockValidationTimeout.RUnlock()
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