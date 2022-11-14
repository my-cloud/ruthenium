// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package network

import (
	"github.com/my-cloud/ruthenium/src/protocol"
	"sync"
)

// Ensure, that ValidatorMock does implement Validator.
// If this is not the case, regenerate this file with moq.
var _ protocol.Validator = &ValidatorMock{}

// ValidatorMock is a mock implementation of Validator.
//
// 	func TestSomethingThatUsesValidator(t *testing.T) {
//
// 		// make and configure a mocked Validator
// 		mockedValidator := &ValidatorMock{
// 			StartValidationFunc: func()  {
// 				panic("mock out the StartValidation method")
// 			},
// 			StopValidationFunc: func()  {
// 				panic("mock out the StopValidation method")
// 			},
// 			ValidateFunc: func()  {
// 				panic("mock out the Validate method")
// 			},
// 		}
//
// 		// use mockedValidator in code that requires Validator
// 		// and then make assertions.
//
// 	}
type ValidatorMock struct {
	// StartValidationFunc mocks the StartValidation method.
	StartValidationFunc func()

	// StopValidationFunc mocks the StopValidation method.
	StopValidationFunc func()

	// ValidateFunc mocks the Validate method.
	ValidateFunc func()

	// calls tracks calls to the methods.
	calls struct {
		// StartValidation holds details about calls to the StartValidation method.
		StartValidation []struct {
		}
		// StopValidation holds details about calls to the StopValidation method.
		StopValidation []struct {
		}
		// Validate holds details about calls to the Validate method.
		Validate []struct {
		}
	}
	lockStartValidation sync.RWMutex
	lockStopValidation  sync.RWMutex
	lockValidate        sync.RWMutex
}

// StartValidation calls StartValidationFunc.
func (mock *ValidatorMock) StartValidation() {
	if mock.StartValidationFunc == nil {
		panic("ValidatorMock.StartValidationFunc: method is nil but Validator.StartValidation was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStartValidation.Lock()
	mock.calls.StartValidation = append(mock.calls.StartValidation, callInfo)
	mock.lockStartValidation.Unlock()
	mock.StartValidationFunc()
}

// StartValidationCalls gets all the calls that were made to StartValidation.
// Check the length with:
//     len(mockedValidator.StartValidationCalls())
func (mock *ValidatorMock) StartValidationCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStartValidation.RLock()
	calls = mock.calls.StartValidation
	mock.lockStartValidation.RUnlock()
	return calls
}

// StopValidation calls StopValidationFunc.
func (mock *ValidatorMock) StopValidation() {
	if mock.StopValidationFunc == nil {
		panic("ValidatorMock.StopValidationFunc: method is nil but Validator.StopValidation was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStopValidation.Lock()
	mock.calls.StopValidation = append(mock.calls.StopValidation, callInfo)
	mock.lockStopValidation.Unlock()
	mock.StopValidationFunc()
}

// StopValidationCalls gets all the calls that were made to StopValidation.
// Check the length with:
//     len(mockedValidator.StopValidationCalls())
func (mock *ValidatorMock) StopValidationCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStopValidation.RLock()
	calls = mock.calls.StopValidation
	mock.lockStopValidation.RUnlock()
	return calls
}

// Validate calls ValidateFunc.
func (mock *ValidatorMock) Validate() {
	if mock.ValidateFunc == nil {
		panic("ValidatorMock.ValidateFunc: method is nil but Validator.Validate was just called")
	}
	callInfo := struct {
	}{}
	mock.lockValidate.Lock()
	mock.calls.Validate = append(mock.calls.Validate, callInfo)
	mock.lockValidate.Unlock()
	mock.ValidateFunc()
}

// ValidateCalls gets all the calls that were made to Validate.
// Check the length with:
//     len(mockedValidator.ValidateCalls())
func (mock *ValidatorMock) ValidateCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockValidate.RLock()
	calls = mock.calls.Validate
	mock.lockValidate.RUnlock()
	return calls
}
