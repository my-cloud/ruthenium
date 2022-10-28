// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package network

import (
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/network"
	"sync"
)

// Ensure, that TransactionsPoolMock does implement TransactionsPool.
// If this is not the case, regenerate this file with moq.
var _ network.TransactionsPool = &TransactionsPoolMock{}

// TransactionsPoolMock is a mock implementation of TransactionsPool.
//
// 	func TestSomethingThatUsesTransactionsPool(t *testing.T) {
//
// 		// make and configure a mocked TransactionsPool
// 		mockedTransactionsPool := &TransactionsPoolMock{
// 			AddTransactionFunc: func(transactionRequest *neighborhood.TransactionRequest, blockchain Blockchain, neighbors []neighborhood.Neighbor)  {
// 				panic("mock out the AddTransaction method")
// 			},
// 			TransactionsFunc: func() []*neighborhood.TransactionResponse {
// 				panic("mock out the Transactions method")
// 			},
// 		}
//
// 		// use mockedTransactionsPool in code that requires TransactionsPool
// 		// and then make assertions.
//
// 	}
type TransactionsPoolMock struct {
	// AddTransactionFunc mocks the AddTransaction method.
	AddTransactionFunc func(transactionRequest *neighborhood.TransactionRequest, blockchain network.Blockchain, neighbors []neighborhood.Neighbor)

	// TransactionsFunc mocks the Transactions method.
	TransactionsFunc func() []*neighborhood.TransactionResponse

	// calls tracks calls to the methods.
	calls struct {
		// AddTransaction holds details about calls to the AddTransaction method.
		AddTransaction []struct {
			// TransactionRequest is the transactionRequest argument value.
			TransactionRequest *neighborhood.TransactionRequest
			// Blockchain is the blockchain argument value.
			Blockchain network.Blockchain
			// Neighbors is the neighbors argument value.
			Neighbors []neighborhood.Neighbor
		}
		// Transactions holds details about calls to the Transactions method.
		Transactions []struct {
		}
	}
	lockAddTransaction sync.RWMutex
	lockTransactions   sync.RWMutex
}

// AddTransaction calls AddTransactionFunc.
func (mock *TransactionsPoolMock) AddTransaction(transactionRequest *neighborhood.TransactionRequest, blockchain network.Blockchain, neighbors []neighborhood.Neighbor) {
	if mock.AddTransactionFunc == nil {
		panic("TransactionsPoolMock.AddTransactionFunc: method is nil but TransactionsPool.AddTransaction was just called")
	}
	callInfo := struct {
		TransactionRequest *neighborhood.TransactionRequest
		Blockchain         network.Blockchain
		Neighbors          []neighborhood.Neighbor
	}{
		TransactionRequest: transactionRequest,
		Blockchain:         blockchain,
		Neighbors:          neighbors,
	}
	mock.lockAddTransaction.Lock()
	mock.calls.AddTransaction = append(mock.calls.AddTransaction, callInfo)
	mock.lockAddTransaction.Unlock()
	mock.AddTransactionFunc(transactionRequest, blockchain, neighbors)
}

// AddTransactionCalls gets all the calls that were made to AddTransaction.
// Check the length with:
//     len(mockedTransactionsPool.AddTransactionCalls())
func (mock *TransactionsPoolMock) AddTransactionCalls() []struct {
	TransactionRequest *neighborhood.TransactionRequest
	Blockchain         network.Blockchain
	Neighbors          []neighborhood.Neighbor
} {
	var calls []struct {
		TransactionRequest *neighborhood.TransactionRequest
		Blockchain         network.Blockchain
		Neighbors          []neighborhood.Neighbor
	}
	mock.lockAddTransaction.RLock()
	calls = mock.calls.AddTransaction
	mock.lockAddTransaction.RUnlock()
	return calls
}

// Transactions calls TransactionsFunc.
func (mock *TransactionsPoolMock) Transactions() []*neighborhood.TransactionResponse {
	if mock.TransactionsFunc == nil {
		panic("TransactionsPoolMock.TransactionsFunc: method is nil but TransactionsPool.Transactions was just called")
	}
	callInfo := struct {
	}{}
	mock.lockTransactions.Lock()
	mock.calls.Transactions = append(mock.calls.Transactions, callInfo)
	mock.lockTransactions.Unlock()
	return mock.TransactionsFunc()
}

// TransactionsCalls gets all the calls that were made to Transactions.
// Check the length with:
//     len(mockedTransactionsPool.TransactionsCalls())
func (mock *TransactionsPoolMock) TransactionsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockTransactions.RLock()
	calls = mock.calls.Transactions
	mock.lockTransactions.RUnlock()
	return calls
}