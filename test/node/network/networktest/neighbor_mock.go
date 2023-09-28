// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package networktest

import (
	"github.com/my-cloud/ruthenium/src/node/network"
	"sync"
)

// Ensure, that NeighborMock does implement Neighbor.
// If this is not the case, regenerate this file with moq.
var _ network.Neighbor = &NeighborMock{}

// NeighborMock is a mock implementation of Neighbor.
//
//	func TestSomethingThatUsesNeighbor(t *testing.T) {
//
//		// make and configure a mocked Neighbor
//		mockedNeighbor := &NeighborMock{
//			AddTransactionFunc: func(request TransactionRequest) error {
//				panic("mock out the AddTransaction method")
//			},
//			GetBlocksFunc: func(startingBlockHeight uint64) ([]byte, error) {
//				panic("mock out the GetBlocks method")
//			},
//			GetFirstBlockTimestampFunc: func() (int64, error) {
//				panic("mock out the GetFirstBlockTimestamp method")
//			},
//			GetTransactionsFunc: func() ([]byte, error) {
//				panic("mock out the GetTransactions method")
//			},
//			GetUtxosFunc: func(address string) ([]*UtxoResponse, error) {
//				panic("mock out the GetUtxos method")
//			},
//			SendTargetsFunc: func(request []TargetRequest) error {
//				panic("mock out the SendTargets method")
//			},
//			TargetFunc: func() string {
//				panic("mock out the Target method")
//			},
//		}
//
//		// use mockedNeighbor in code that requires Neighbor
//		// and then make assertions.
//
//	}
type NeighborMock struct {
	// AddTransactionFunc mocks the AddTransaction method.
	AddTransactionFunc func(request network.TransactionRequest) error

	// GetBlocksFunc mocks the GetBlocks method.
	GetBlocksFunc func(startingBlockHeight uint64) ([]byte, error)

	// GetFirstBlockTimestampFunc mocks the GetFirstBlockTimestamp method.
	GetFirstBlockTimestampFunc func() (int64, error)

	// GetTransactionsFunc mocks the GetTransactions method.
	GetTransactionsFunc func() ([]byte, error)

	// GetUtxosFunc mocks the GetUtxos method.
	GetUtxosFunc func(address string) ([]*network.UtxoResponse, error)

	// SendTargetsFunc mocks the SendTargets method.
	SendTargetsFunc func(request []network.TargetRequest) error

	// TargetFunc mocks the Target method.
	TargetFunc func() string

	// calls tracks calls to the methods.
	calls struct {
		// AddTransaction holds details about calls to the AddTransaction method.
		AddTransaction []struct {
			// Request is the request argument value.
			Request network.TransactionRequest
		}
		// GetBlocks holds details about calls to the GetBlocks method.
		GetBlocks []struct {
			// StartingBlockHeight is the startingBlockHeight argument value.
			StartingBlockHeight uint64
		}
		// GetFirstBlockTimestamp holds details about calls to the GetFirstBlockTimestamp method.
		GetFirstBlockTimestamp []struct {
		}
		// GetTransactions holds details about calls to the GetTransactions method.
		GetTransactions []struct {
		}
		// GetUtxos holds details about calls to the GetUtxos method.
		GetUtxos []struct {
			// Address is the address argument value.
			Address string
		}
		// SendTargets holds details about calls to the SendTargets method.
		SendTargets []struct {
			// Request is the request argument value.
			Request []network.TargetRequest
		}
		// Target holds details about calls to the Target method.
		Target []struct {
		}
	}
	lockAddTransaction         sync.RWMutex
	lockGetBlocks              sync.RWMutex
	lockGetFirstBlockTimestamp sync.RWMutex
	lockGetTransactions        sync.RWMutex
	lockGetUtxos               sync.RWMutex
	lockSendTargets            sync.RWMutex
	lockTarget                 sync.RWMutex
}

// AddTransaction calls AddTransactionFunc.
func (mock *NeighborMock) AddTransaction(request network.TransactionRequest) error {
	if mock.AddTransactionFunc == nil {
		panic("NeighborMock.AddTransactionFunc: method is nil but Neighbor.AddTransaction was just called")
	}
	callInfo := struct {
		Request network.TransactionRequest
	}{
		Request: request,
	}
	mock.lockAddTransaction.Lock()
	mock.calls.AddTransaction = append(mock.calls.AddTransaction, callInfo)
	mock.lockAddTransaction.Unlock()
	return mock.AddTransactionFunc(request)
}

// AddTransactionCalls gets all the calls that were made to AddTransaction.
// Check the length with:
//
//	len(mockedNeighbor.AddTransactionCalls())
func (mock *NeighborMock) AddTransactionCalls() []struct {
	Request network.TransactionRequest
} {
	var calls []struct {
		Request network.TransactionRequest
	}
	mock.lockAddTransaction.RLock()
	calls = mock.calls.AddTransaction
	mock.lockAddTransaction.RUnlock()
	return calls
}

// GetBlocks calls GetBlocksFunc.
func (mock *NeighborMock) GetBlocks(startingBlockHeight uint64) ([]byte, error) {
	if mock.GetBlocksFunc == nil {
		panic("NeighborMock.GetBlocksFunc: method is nil but Neighbor.GetBlocks was just called")
	}
	callInfo := struct {
		StartingBlockHeight uint64
	}{
		StartingBlockHeight: startingBlockHeight,
	}
	mock.lockGetBlocks.Lock()
	mock.calls.GetBlocks = append(mock.calls.GetBlocks, callInfo)
	mock.lockGetBlocks.Unlock()
	return mock.GetBlocksFunc(startingBlockHeight)
}

// GetBlocksCalls gets all the calls that were made to GetBlocks.
// Check the length with:
//
//	len(mockedNeighbor.GetBlocksCalls())
func (mock *NeighborMock) GetBlocksCalls() []struct {
	StartingBlockHeight uint64
} {
	var calls []struct {
		StartingBlockHeight uint64
	}
	mock.lockGetBlocks.RLock()
	calls = mock.calls.GetBlocks
	mock.lockGetBlocks.RUnlock()
	return calls
}

// GetFirstBlockTimestamp calls GetFirstBlockTimestampFunc.
func (mock *NeighborMock) GetFirstBlockTimestamp() (int64, error) {
	if mock.GetFirstBlockTimestampFunc == nil {
		panic("NeighborMock.GetFirstBlockTimestampFunc: method is nil but Neighbor.GetFirstBlockTimestamp was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetFirstBlockTimestamp.Lock()
	mock.calls.GetFirstBlockTimestamp = append(mock.calls.GetFirstBlockTimestamp, callInfo)
	mock.lockGetFirstBlockTimestamp.Unlock()
	return mock.GetFirstBlockTimestampFunc()
}

// GetFirstBlockTimestampCalls gets all the calls that were made to GetFirstBlockTimestamp.
// Check the length with:
//
//	len(mockedNeighbor.GetFirstBlockTimestampCalls())
func (mock *NeighborMock) GetFirstBlockTimestampCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetFirstBlockTimestamp.RLock()
	calls = mock.calls.GetFirstBlockTimestamp
	mock.lockGetFirstBlockTimestamp.RUnlock()
	return calls
}

// GetTransactions calls GetTransactionsFunc.
func (mock *NeighborMock) GetTransactions() ([]byte, error) {
	if mock.GetTransactionsFunc == nil {
		panic("NeighborMock.GetTransactionsFunc: method is nil but Neighbor.GetTransactions was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetTransactions.Lock()
	mock.calls.GetTransactions = append(mock.calls.GetTransactions, callInfo)
	mock.lockGetTransactions.Unlock()
	return mock.GetTransactionsFunc()
}

// GetTransactionsCalls gets all the calls that were made to GetTransactions.
// Check the length with:
//
//	len(mockedNeighbor.GetTransactionsCalls())
func (mock *NeighborMock) GetTransactionsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetTransactions.RLock()
	calls = mock.calls.GetTransactions
	mock.lockGetTransactions.RUnlock()
	return calls
}

// GetUtxos calls GetUtxosFunc.
func (mock *NeighborMock) GetUtxos(address string) ([]*network.UtxoResponse, error) {
	if mock.GetUtxosFunc == nil {
		panic("NeighborMock.GetUtxosFunc: method is nil but Neighbor.GetUtxos was just called")
	}
	callInfo := struct {
		Address string
	}{
		Address: address,
	}
	mock.lockGetUtxos.Lock()
	mock.calls.GetUtxos = append(mock.calls.GetUtxos, callInfo)
	mock.lockGetUtxos.Unlock()
	return mock.GetUtxosFunc(address)
}

// GetUtxosCalls gets all the calls that were made to GetUtxos.
// Check the length with:
//
//	len(mockedNeighbor.GetUtxosCalls())
func (mock *NeighborMock) GetUtxosCalls() []struct {
	Address string
} {
	var calls []struct {
		Address string
	}
	mock.lockGetUtxos.RLock()
	calls = mock.calls.GetUtxos
	mock.lockGetUtxos.RUnlock()
	return calls
}

// SendTargets calls SendTargetsFunc.
func (mock *NeighborMock) SendTargets(request []network.TargetRequest) error {
	if mock.SendTargetsFunc == nil {
		panic("NeighborMock.SendTargetsFunc: method is nil but Neighbor.SendTargets was just called")
	}
	callInfo := struct {
		Request []network.TargetRequest
	}{
		Request: request,
	}
	mock.lockSendTargets.Lock()
	mock.calls.SendTargets = append(mock.calls.SendTargets, callInfo)
	mock.lockSendTargets.Unlock()
	return mock.SendTargetsFunc(request)
}

// SendTargetsCalls gets all the calls that were made to SendTargets.
// Check the length with:
//
//	len(mockedNeighbor.SendTargetsCalls())
func (mock *NeighborMock) SendTargetsCalls() []struct {
	Request []network.TargetRequest
} {
	var calls []struct {
		Request []network.TargetRequest
	}
	mock.lockSendTargets.RLock()
	calls = mock.calls.SendTargets
	mock.lockSendTargets.RUnlock()
	return calls
}

// Target calls TargetFunc.
func (mock *NeighborMock) Target() string {
	if mock.TargetFunc == nil {
		panic("NeighborMock.TargetFunc: method is nil but Neighbor.Target was just called")
	}
	callInfo := struct {
	}{}
	mock.lockTarget.Lock()
	mock.calls.Target = append(mock.calls.Target, callInfo)
	mock.lockTarget.Unlock()
	return mock.TargetFunc()
}

// TargetCalls gets all the calls that were made to Target.
// Check the length with:
//
//	len(mockedNeighbor.TargetCalls())
func (mock *NeighborMock) TargetCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockTarget.RLock()
	calls = mock.calls.Target
	mock.lockTarget.RUnlock()
	return calls
}
