// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package protocoltest

import (
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"sync"
)

// Ensure, that BlockchainMock does implement Blockchain.
// If this is not the case, regenerate this file with moq.
var _ protocol.Blockchain = &BlockchainMock{}

// BlockchainMock is a mock implementation of Blockchain.
//
//	func TestSomethingThatUsesBlockchain(t *testing.T) {
//
//		// make and configure a mocked Blockchain
//		mockedBlockchain := &BlockchainMock{
//			AddBlockFunc: func(timestamp int64, transactions []*network.TransactionResponse, registeredAddresses []string)  {
//				panic("mock out the AddBlock method")
//			},
//			BlocksFunc: func() []*network.BlockResponse {
//				panic("mock out the Blocks method")
//			},
//			CalculateTotalAmountFunc: func(currentTimestamp int64, blockchainAddress string) uint64 {
//				panic("mock out the CalculateTotalAmount method")
//			},
//			CopyFunc: func() Blockchain {
//				panic("mock out the Copy method")
//			},
//			LastBlocksFunc: func(startingBlockNonce int) []*network.BlockResponse {
//				panic("mock out the LastBlocks method")
//			},
//		}
//
//		// use mockedBlockchain in code that requires Blockchain
//		// and then make assertions.
//
//	}
type BlockchainMock struct {
	// AddBlockFunc mocks the AddBlock method.
	AddBlockFunc func(timestamp int64, transactions []*network.TransactionResponse, registeredAddresses []string)

	// BlocksFunc mocks the Blocks method.
	BlocksFunc func() []*network.BlockResponse

	// CalculateTotalAmountFunc mocks the CalculateTotalAmount method.
	CalculateTotalAmountFunc func(currentTimestamp int64, blockchainAddress string) uint64

	// CopyFunc mocks the Copy method.
	CopyFunc func() protocol.Blockchain

	// LastBlocksFunc mocks the LastBlocks method.
	LastBlocksFunc func(startingBlockNonce int) []*network.BlockResponse

	// calls tracks calls to the methods.
	calls struct {
		// AddBlock holds details about calls to the AddBlock method.
		AddBlock []struct {
			// Timestamp is the timestamp argument value.
			Timestamp int64
			// Transactions is the transactions argument value.
			Transactions []*network.TransactionResponse
			// RegisteredAddresses is the registeredAddresses argument value.
			RegisteredAddresses []string
		}
		// Blocks holds details about calls to the Blocks method.
		Blocks []struct {
		}
		// CalculateTotalAmount holds details about calls to the CalculateTotalAmount method.
		CalculateTotalAmount []struct {
			// CurrentTimestamp is the currentTimestamp argument value.
			CurrentTimestamp int64
			// BlockchainAddress is the blockchainAddress argument value.
			BlockchainAddress string
		}
		// Copy holds details about calls to the Copy method.
		Copy []struct {
		}
		// LastBlocks holds details about calls to the LastBlocks method.
		LastBlocks []struct {
			// StartingBlockNonce is the startingBlockNonce argument value.
			StartingBlockNonce int
		}
	}
	lockAddBlock             sync.RWMutex
	lockBlocks               sync.RWMutex
	lockCalculateTotalAmount sync.RWMutex
	lockCopy                 sync.RWMutex
	lockLastBlocks           sync.RWMutex
}

// AddBlock calls AddBlockFunc.
func (mock *BlockchainMock) AddBlock(timestamp int64, transactions []*network.TransactionResponse, registeredAddresses []string) {
	if mock.AddBlockFunc == nil {
		panic("BlockchainMock.AddBlockFunc: method is nil but Blockchain.AddBlock was just called")
	}
	callInfo := struct {
		Timestamp           int64
		Transactions        []*network.TransactionResponse
		RegisteredAddresses []string
	}{
		Timestamp:           timestamp,
		Transactions:        transactions,
		RegisteredAddresses: registeredAddresses,
	}
	mock.lockAddBlock.Lock()
	mock.calls.AddBlock = append(mock.calls.AddBlock, callInfo)
	mock.lockAddBlock.Unlock()
	mock.AddBlockFunc(timestamp, transactions, registeredAddresses)
}

// AddBlockCalls gets all the calls that were made to AddBlock.
// Check the length with:
//
//	len(mockedBlockchain.AddBlockCalls())
func (mock *BlockchainMock) AddBlockCalls() []struct {
	Timestamp           int64
	Transactions        []*network.TransactionResponse
	RegisteredAddresses []string
} {
	var calls []struct {
		Timestamp           int64
		Transactions        []*network.TransactionResponse
		RegisteredAddresses []string
	}
	mock.lockAddBlock.RLock()
	calls = mock.calls.AddBlock
	mock.lockAddBlock.RUnlock()
	return calls
}

// Blocks calls BlocksFunc.
func (mock *BlockchainMock) Blocks() []*network.BlockResponse {
	if mock.BlocksFunc == nil {
		panic("BlockchainMock.BlocksFunc: method is nil but Blockchain.Blocks was just called")
	}
	callInfo := struct {
	}{}
	mock.lockBlocks.Lock()
	mock.calls.Blocks = append(mock.calls.Blocks, callInfo)
	mock.lockBlocks.Unlock()
	return mock.BlocksFunc()
}

// BlocksCalls gets all the calls that were made to Blocks.
// Check the length with:
//
//	len(mockedBlockchain.BlocksCalls())
func (mock *BlockchainMock) BlocksCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockBlocks.RLock()
	calls = mock.calls.Blocks
	mock.lockBlocks.RUnlock()
	return calls
}

// CalculateTotalAmount calls CalculateTotalAmountFunc.
func (mock *BlockchainMock) CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64 {
	if mock.CalculateTotalAmountFunc == nil {
		panic("BlockchainMock.CalculateTotalAmountFunc: method is nil but Blockchain.CalculateTotalAmount was just called")
	}
	callInfo := struct {
		CurrentTimestamp  int64
		BlockchainAddress string
	}{
		CurrentTimestamp:  currentTimestamp,
		BlockchainAddress: blockchainAddress,
	}
	mock.lockCalculateTotalAmount.Lock()
	mock.calls.CalculateTotalAmount = append(mock.calls.CalculateTotalAmount, callInfo)
	mock.lockCalculateTotalAmount.Unlock()
	return mock.CalculateTotalAmountFunc(currentTimestamp, blockchainAddress)
}

// CalculateTotalAmountCalls gets all the calls that were made to CalculateTotalAmount.
// Check the length with:
//
//	len(mockedBlockchain.CalculateTotalAmountCalls())
func (mock *BlockchainMock) CalculateTotalAmountCalls() []struct {
	CurrentTimestamp  int64
	BlockchainAddress string
} {
	var calls []struct {
		CurrentTimestamp  int64
		BlockchainAddress string
	}
	mock.lockCalculateTotalAmount.RLock()
	calls = mock.calls.CalculateTotalAmount
	mock.lockCalculateTotalAmount.RUnlock()
	return calls
}

// Copy calls CopyFunc.
func (mock *BlockchainMock) Copy() protocol.Blockchain {
	if mock.CopyFunc == nil {
		panic("BlockchainMock.CopyFunc: method is nil but Blockchain.Copy was just called")
	}
	callInfo := struct {
	}{}
	mock.lockCopy.Lock()
	mock.calls.Copy = append(mock.calls.Copy, callInfo)
	mock.lockCopy.Unlock()
	return mock.CopyFunc()
}

// CopyCalls gets all the calls that were made to Copy.
// Check the length with:
//
//	len(mockedBlockchain.CopyCalls())
func (mock *BlockchainMock) CopyCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockCopy.RLock()
	calls = mock.calls.Copy
	mock.lockCopy.RUnlock()
	return calls
}

// LastBlocks calls LastBlocksFunc.
func (mock *BlockchainMock) LastBlocks(startingBlockNonce int) []*network.BlockResponse {
	if mock.LastBlocksFunc == nil {
		panic("BlockchainMock.LastBlocksFunc: method is nil but Blockchain.LastBlocks was just called")
	}
	callInfo := struct {
		StartingBlockNonce int
	}{
		StartingBlockNonce: startingBlockNonce,
	}
	mock.lockLastBlocks.Lock()
	mock.calls.LastBlocks = append(mock.calls.LastBlocks, callInfo)
	mock.lockLastBlocks.Unlock()
	return mock.LastBlocksFunc(startingBlockNonce)
}

// LastBlocksCalls gets all the calls that were made to LastBlocks.
// Check the length with:
//
//	len(mockedBlockchain.LastBlocksCalls())
func (mock *BlockchainMock) LastBlocksCalls() []struct {
	StartingBlockNonce int
} {
	var calls []struct {
		StartingBlockNonce int
	}
	mock.lockLastBlocks.RLock()
	calls = mock.calls.LastBlocks
	mock.lockLastBlocks.RUnlock()
	return calls
}
