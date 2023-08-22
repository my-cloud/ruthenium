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
//			AddBlockFunc: func(timestamp int64, transactions []*network.TransactionResponse, newRegisteredAddresses []string) error {
//				panic("mock out the AddBlock method")
//			},
//			BlockFunc: func(blockHeight uint64) *network.BlockResponse {
//				panic("mock out the Block method")
//			},
//			BlocksFunc: func() []*network.BlockResponse {
//				panic("mock out the Blocks method")
//			},
//			CopyFunc: func() Blockchain {
//				panic("mock out the Copy method")
//			},
//			FindFeeFunc: func(transaction *network.TransactionResponse, timestamp int64) (uint64, error) {
//				panic("mock out the FindFee method")
//			},
//			LastBlocksFunc: func(startingBlockHeight uint64) []*network.BlockResponse {
//				panic("mock out the LastBlocks method")
//			},
//			UtxosByAddressFunc: func(address string) []*network.UtxoResponse {
//				panic("mock out the UtxosByAddress method")
//			},
//		}
//
//		// use mockedBlockchain in code that requires Blockchain
//		// and then make assertions.
//
//	}
type BlockchainMock struct {
	// AddBlockFunc mocks the AddBlock method.
	AddBlockFunc func(timestamp int64, transactions []*network.TransactionResponse, newRegisteredAddresses []string) error

	// BlockFunc mocks the Block method.
	BlockFunc func(blockHeight uint64) *network.BlockResponse

	// BlocksFunc mocks the Blocks method.
	BlocksFunc func() []*network.BlockResponse

	// CopyFunc mocks the Copy method.
	CopyFunc func() protocol.Blockchain

	// FindFeeFunc mocks the FindFee method.
	FindFeeFunc func(transaction *network.TransactionResponse, timestamp int64) (uint64, error)

	// LastBlocksFunc mocks the LastBlocks method.
	LastBlocksFunc func(startingBlockHeight uint64) []*network.BlockResponse

	// UtxosByAddressFunc mocks the UtxosByAddress method.
	UtxosByAddressFunc func(address string) []*network.UtxoResponse

	// calls tracks calls to the methods.
	calls struct {
		// AddBlock holds details about calls to the AddBlock method.
		AddBlock []struct {
			// Timestamp is the timestamp argument value.
			Timestamp int64
			// Transactions is the transactions argument value.
			Transactions []*network.TransactionResponse
			// NewRegisteredAddresses is the newRegisteredAddresses argument value.
			NewRegisteredAddresses []string
		}
		// Block holds details about calls to the Block method.
		Block []struct {
			// BlockHeight is the blockHeight argument value.
			BlockHeight uint64
		}
		// Blocks holds details about calls to the Blocks method.
		Blocks []struct {
		}
		// Copy holds details about calls to the Copy method.
		Copy []struct {
		}
		// FindFee holds details about calls to the FindFee method.
		FindFee []struct {
			// Transaction is the transaction argument value.
			Transaction *network.TransactionResponse
			// Timestamp is the timestamp argument value.
			Timestamp int64
		}
		// LastBlocks holds details about calls to the LastBlocks method.
		LastBlocks []struct {
			// StartingBlockHeight is the startingBlockHeight argument value.
			StartingBlockHeight uint64
		}
		// UtxosByAddress holds details about calls to the UtxosByAddress method.
		UtxosByAddress []struct {
			// Address is the address argument value.
			Address string
		}
	}
	lockAddBlock       sync.RWMutex
	lockBlock          sync.RWMutex
	lockBlocks         sync.RWMutex
	lockCopy           sync.RWMutex
	lockFindFee        sync.RWMutex
	lockLastBlocks     sync.RWMutex
	lockUtxosByAddress sync.RWMutex
}

// AddBlock calls AddBlockFunc.
func (mock *BlockchainMock) AddBlock(timestamp int64, transactions []*network.TransactionResponse, newRegisteredAddresses []string) error {
	if mock.AddBlockFunc == nil {
		panic("BlockchainMock.AddBlockFunc: method is nil but Blockchain.AddBlock was just called")
	}
	callInfo := struct {
		Timestamp              int64
		Transactions           []*network.TransactionResponse
		NewRegisteredAddresses []string
	}{
		Timestamp:              timestamp,
		Transactions:           transactions,
		NewRegisteredAddresses: newRegisteredAddresses,
	}
	mock.lockAddBlock.Lock()
	mock.calls.AddBlock = append(mock.calls.AddBlock, callInfo)
	mock.lockAddBlock.Unlock()
	return mock.AddBlockFunc(timestamp, transactions, newRegisteredAddresses)
}

// AddBlockCalls gets all the calls that were made to AddBlock.
// Check the length with:
//
//	len(mockedBlockchain.AddBlockCalls())
func (mock *BlockchainMock) AddBlockCalls() []struct {
	Timestamp              int64
	Transactions           []*network.TransactionResponse
	NewRegisteredAddresses []string
} {
	var calls []struct {
		Timestamp              int64
		Transactions           []*network.TransactionResponse
		NewRegisteredAddresses []string
	}
	mock.lockAddBlock.RLock()
	calls = mock.calls.AddBlock
	mock.lockAddBlock.RUnlock()
	return calls
}

// Block calls BlockFunc.
func (mock *BlockchainMock) Block(blockHeight uint64) *network.BlockResponse {
	if mock.BlockFunc == nil {
		panic("BlockchainMock.BlockFunc: method is nil but Blockchain.Block was just called")
	}
	callInfo := struct {
		BlockHeight uint64
	}{
		BlockHeight: blockHeight,
	}
	mock.lockBlock.Lock()
	mock.calls.Block = append(mock.calls.Block, callInfo)
	mock.lockBlock.Unlock()
	return mock.BlockFunc(blockHeight)
}

// BlockCalls gets all the calls that were made to Block.
// Check the length with:
//
//	len(mockedBlockchain.BlockCalls())
func (mock *BlockchainMock) BlockCalls() []struct {
	BlockHeight uint64
} {
	var calls []struct {
		BlockHeight uint64
	}
	mock.lockBlock.RLock()
	calls = mock.calls.Block
	mock.lockBlock.RUnlock()
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

// FindFee calls FindFeeFunc.
func (mock *BlockchainMock) FindFee(transaction *network.TransactionResponse, timestamp int64) (uint64, error) {
	if mock.FindFeeFunc == nil {
		panic("BlockchainMock.FindFeeFunc: method is nil but Blockchain.FindFee was just called")
	}
	callInfo := struct {
		Transaction *network.TransactionResponse
		Timestamp   int64
	}{
		Transaction: transaction,
		Timestamp:   timestamp,
	}
	mock.lockFindFee.Lock()
	mock.calls.FindFee = append(mock.calls.FindFee, callInfo)
	mock.lockFindFee.Unlock()
	return mock.FindFeeFunc(transaction, timestamp)
}

// FindFeeCalls gets all the calls that were made to FindFee.
// Check the length with:
//
//	len(mockedBlockchain.FindFeeCalls())
func (mock *BlockchainMock) FindFeeCalls() []struct {
	Transaction *network.TransactionResponse
	Timestamp   int64
} {
	var calls []struct {
		Transaction *network.TransactionResponse
		Timestamp   int64
	}
	mock.lockFindFee.RLock()
	calls = mock.calls.FindFee
	mock.lockFindFee.RUnlock()
	return calls
}

// LastBlocks calls LastBlocksFunc.
func (mock *BlockchainMock) LastBlocks(startingBlockHeight uint64) []*network.BlockResponse {
	if mock.LastBlocksFunc == nil {
		panic("BlockchainMock.LastBlocksFunc: method is nil but Blockchain.LastBlocks was just called")
	}
	callInfo := struct {
		StartingBlockHeight uint64
	}{
		StartingBlockHeight: startingBlockHeight,
	}
	mock.lockLastBlocks.Lock()
	mock.calls.LastBlocks = append(mock.calls.LastBlocks, callInfo)
	mock.lockLastBlocks.Unlock()
	return mock.LastBlocksFunc(startingBlockHeight)
}

// LastBlocksCalls gets all the calls that were made to LastBlocks.
// Check the length with:
//
//	len(mockedBlockchain.LastBlocksCalls())
func (mock *BlockchainMock) LastBlocksCalls() []struct {
	StartingBlockHeight uint64
} {
	var calls []struct {
		StartingBlockHeight uint64
	}
	mock.lockLastBlocks.RLock()
	calls = mock.calls.LastBlocks
	mock.lockLastBlocks.RUnlock()
	return calls
}

// UtxosByAddress calls UtxosByAddressFunc.
func (mock *BlockchainMock) UtxosByAddress(address string) []*network.UtxoResponse {
	if mock.UtxosByAddressFunc == nil {
		panic("BlockchainMock.UtxosByAddressFunc: method is nil but Blockchain.UtxosByAddress was just called")
	}
	callInfo := struct {
		Address string
	}{
		Address: address,
	}
	mock.lockUtxosByAddress.Lock()
	mock.calls.UtxosByAddress = append(mock.calls.UtxosByAddress, callInfo)
	mock.lockUtxosByAddress.Unlock()
	return mock.UtxosByAddressFunc(address)
}

// UtxosByAddressCalls gets all the calls that were made to UtxosByAddress.
// Check the length with:
//
//	len(mockedBlockchain.UtxosByAddressCalls())
func (mock *BlockchainMock) UtxosByAddressCalls() []struct {
	Address string
} {
	var calls []struct {
		Address string
	}
	mock.lockUtxosByAddress.RLock()
	calls = mock.calls.UtxosByAddress
	mock.lockUtxosByAddress.RUnlock()
	return calls
}
