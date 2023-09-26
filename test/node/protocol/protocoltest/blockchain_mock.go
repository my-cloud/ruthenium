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
//			AddBlockFunc: func(timestamp int64, transactions []byte, newRegisteredAddresses []string) error {
//				panic("mock out the AddBlock method")
//			},
//			BlocksFunc: func(startingBlockHeight uint64) []byte {
//				panic("mock out the Blocks method")
//			},
//			CopyFunc: func() Blockchain {
//				panic("mock out the Copy method")
//			},
//			FindFeeFunc: func(inputs []*network.InputResponse, outputs []*network.OutputResponse, timestamp int64) (uint64, error) {
//				panic("mock out the FindFee method")
//			},
//			FirstBlockTimestampFunc: func() int64 {
//				panic("mock out the FirstBlockTimestamp method")
//			},
//			LastBlockTimestampFunc: func() int64 {
//				panic("mock out the LastBlockTimestamp method")
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
	AddBlockFunc func(timestamp int64, transactions []byte, newRegisteredAddresses []string) error

	// BlocksFunc mocks the Blocks method.
	BlocksFunc func(startingBlockHeight uint64) []byte

	// CopyFunc mocks the Copy method.
	CopyFunc func() protocol.Blockchain

	// FindFeeFunc mocks the FindFee method.
	FindFeeFunc func(inputs []*network.InputResponse, outputs []*network.OutputResponse, timestamp int64) (uint64, error)

	// FirstBlockTimestampFunc mocks the FirstBlockTimestamp method.
	FirstBlockTimestampFunc func() int64

	// LastBlockTimestampFunc mocks the LastBlockTimestamp method.
	LastBlockTimestampFunc func() int64

	// UtxosByAddressFunc mocks the UtxosByAddress method.
	UtxosByAddressFunc func(address string) []*network.UtxoResponse

	// calls tracks calls to the methods.
	calls struct {
		// AddBlock holds details about calls to the AddBlock method.
		AddBlock []struct {
			// Timestamp is the timestamp argument value.
			Timestamp int64
			// Transactions is the transactions argument value.
			Transactions []byte
			// NewRegisteredAddresses is the newRegisteredAddresses argument value.
			NewRegisteredAddresses []string
		}
		// Blocks holds details about calls to the Blocks method.
		Blocks []struct {
			// StartingBlockHeight is the startingBlockHeight argument value.
			StartingBlockHeight uint64
		}
		// Copy holds details about calls to the Copy method.
		Copy []struct {
		}
		// FindFee holds details about calls to the FindFee method.
		FindFee []struct {
			// Inputs is the inputs argument value.
			Inputs []*network.InputResponse
			// Outputs is the outputs argument value.
			Outputs []*network.OutputResponse
			// Timestamp is the timestamp argument value.
			Timestamp int64
		}
		// FirstBlockTimestamp holds details about calls to the FirstBlockTimestamp method.
		FirstBlockTimestamp []struct {
		}
		// LastBlockTimestamp holds details about calls to the LastBlockTimestamp method.
		LastBlockTimestamp []struct {
		}
		// UtxosByAddress holds details about calls to the UtxosByAddress method.
		UtxosByAddress []struct {
			// Address is the address argument value.
			Address string
		}
	}
	lockAddBlock            sync.RWMutex
	lockBlocks              sync.RWMutex
	lockCopy                sync.RWMutex
	lockFindFee             sync.RWMutex
	lockFirstBlockTimestamp sync.RWMutex
	lockLastBlockTimestamp  sync.RWMutex
	lockUtxosByAddress      sync.RWMutex
}

// AddBlock calls AddBlockFunc.
func (mock *BlockchainMock) AddBlock(timestamp int64, transactions []byte, newRegisteredAddresses []string) error {
	if mock.AddBlockFunc == nil {
		panic("BlockchainMock.AddBlockFunc: method is nil but Blockchain.AddBlock was just called")
	}
	callInfo := struct {
		Timestamp              int64
		Transactions           []byte
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
	Transactions           []byte
	NewRegisteredAddresses []string
} {
	var calls []struct {
		Timestamp              int64
		Transactions           []byte
		NewRegisteredAddresses []string
	}
	mock.lockAddBlock.RLock()
	calls = mock.calls.AddBlock
	mock.lockAddBlock.RUnlock()
	return calls
}

// Blocks calls BlocksFunc.
func (mock *BlockchainMock) Blocks(startingBlockHeight uint64) []byte {
	if mock.BlocksFunc == nil {
		panic("BlockchainMock.BlocksFunc: method is nil but Blockchain.Blocks was just called")
	}
	callInfo := struct {
		StartingBlockHeight uint64
	}{
		StartingBlockHeight: startingBlockHeight,
	}
	mock.lockBlocks.Lock()
	mock.calls.Blocks = append(mock.calls.Blocks, callInfo)
	mock.lockBlocks.Unlock()
	return mock.BlocksFunc(startingBlockHeight)
}

// BlocksCalls gets all the calls that were made to Blocks.
// Check the length with:
//
//	len(mockedBlockchain.BlocksCalls())
func (mock *BlockchainMock) BlocksCalls() []struct {
	StartingBlockHeight uint64
} {
	var calls []struct {
		StartingBlockHeight uint64
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
func (mock *BlockchainMock) FindFee(inputs []*network.InputResponse, outputs []*network.OutputResponse, timestamp int64) (uint64, error) {
	if mock.FindFeeFunc == nil {
		panic("BlockchainMock.FindFeeFunc: method is nil but Blockchain.FindFee was just called")
	}
	callInfo := struct {
		Inputs    []*network.InputResponse
		Outputs   []*network.OutputResponse
		Timestamp int64
	}{
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: timestamp,
	}
	mock.lockFindFee.Lock()
	mock.calls.FindFee = append(mock.calls.FindFee, callInfo)
	mock.lockFindFee.Unlock()
	return mock.FindFeeFunc(inputs, outputs, timestamp)
}

// FindFeeCalls gets all the calls that were made to FindFee.
// Check the length with:
//
//	len(mockedBlockchain.FindFeeCalls())
func (mock *BlockchainMock) FindFeeCalls() []struct {
	Inputs    []*network.InputResponse
	Outputs   []*network.OutputResponse
	Timestamp int64
} {
	var calls []struct {
		Inputs    []*network.InputResponse
		Outputs   []*network.OutputResponse
		Timestamp int64
	}
	mock.lockFindFee.RLock()
	calls = mock.calls.FindFee
	mock.lockFindFee.RUnlock()
	return calls
}

// FirstBlockTimestamp calls FirstBlockTimestampFunc.
func (mock *BlockchainMock) FirstBlockTimestamp() int64 {
	if mock.FirstBlockTimestampFunc == nil {
		panic("BlockchainMock.FirstBlockTimestampFunc: method is nil but Blockchain.FirstBlockTimestamp was just called")
	}
	callInfo := struct {
	}{}
	mock.lockFirstBlockTimestamp.Lock()
	mock.calls.FirstBlockTimestamp = append(mock.calls.FirstBlockTimestamp, callInfo)
	mock.lockFirstBlockTimestamp.Unlock()
	return mock.FirstBlockTimestampFunc()
}

// FirstBlockTimestampCalls gets all the calls that were made to FirstBlockTimestamp.
// Check the length with:
//
//	len(mockedBlockchain.FirstBlockTimestampCalls())
func (mock *BlockchainMock) FirstBlockTimestampCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockFirstBlockTimestamp.RLock()
	calls = mock.calls.FirstBlockTimestamp
	mock.lockFirstBlockTimestamp.RUnlock()
	return calls
}

// LastBlockTimestamp calls LastBlockTimestampFunc.
func (mock *BlockchainMock) LastBlockTimestamp() int64 {
	if mock.LastBlockTimestampFunc == nil {
		panic("BlockchainMock.LastBlockTimestampFunc: method is nil but Blockchain.LastBlockTimestamp was just called")
	}
	callInfo := struct {
	}{}
	mock.lockLastBlockTimestamp.Lock()
	mock.calls.LastBlockTimestamp = append(mock.calls.LastBlockTimestamp, callInfo)
	mock.lockLastBlockTimestamp.Unlock()
	return mock.LastBlockTimestampFunc()
}

// LastBlockTimestampCalls gets all the calls that were made to LastBlockTimestamp.
// Check the length with:
//
//	len(mockedBlockchain.LastBlockTimestampCalls())
func (mock *BlockchainMock) LastBlockTimestampCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockLastBlockTimestamp.RLock()
	calls = mock.calls.LastBlockTimestamp
	mock.lockLastBlockTimestamp.RUnlock()
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
