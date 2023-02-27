package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocoltest"
	"strings"
	"testing"
)

const (
	blockchainReplacedMessage = "verification done: blockchain replaced"
	blockchainKeptMessage     = "verification done: blockchain kept"
)

func Test_AddBlock_ValidParameters_NoErrorReturned(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, 0, registry, 1, synchronizer, logger)

	// Act
	err := blockchain.AddBlock(0, nil, nil)

	// Assert
	test.Assert(t, err == nil, "error is returned whereas it should not")
}

func Test_Blocks_ValidParameters_NoErrorLogged(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, 0, registry, 1, synchronizer, logger)

	// Act
	blocks := blockchain.Blocks()

	// Assert
	test.Assert(t, len(blocks) == 1, "blocks don't contain a single block")
}

func Test_CalculateTotalAmount_InitialValidator_ReturnsGenesisAmount(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	genesisTransaction := validation.NewRewardTransaction("", 0, 10)
	blockchain := verification.NewBlockchain(genesisTransaction, 0, registry, 1, synchronizer, logger)

	// Act
	amount := blockchain.CalculateTotalAmount(1, genesisTransaction.RecipientAddress)

	// Assert
	test.Assert(t, amount == genesisTransaction.Value, "calculated amount is not the genesis amount whereas it should be")
}

func Test_Update_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, 0, registry, 1, synchronizer, logger)
	_ = blockchain.AddBlock(1, nil, nil)
	_ = blockchain.AddBlock(2, nil, nil)
	_ = blockchain.AddBlock(3, nil, nil)
	_ = blockchain.AddBlock(4, nil, nil)
	neighborMock.GetLastBlocksFunc = func(uint64) ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewRewardedBlockResponse(blockchain.LastBlocks(2)[0].PreviousHash, 2)
		block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
		hash1, _ := block1.Hash()
		blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, 3)
		block2, _ := verification.NewBlockFromResponse(blockResponse2, nil)
		hash2, _ := block2.Hash()
		blockResponse3 := protocoltest.NewRewardedBlockResponse(hash2, 4)
		block3, _ := verification.NewBlockFromResponse(blockResponse3, nil)
		hash3, _ := block3.Hash()
		blockResponse4 := protocoltest.NewRewardedBlockResponse(hash3, 5)
		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3, blockResponse4}, nil
	}

	// Act
	blockchain.Update(5)

	// Assert
	var isReplaced bool
	for _, call := range logger.DebugCalls() {
		if call.Msg == blockchainReplacedMessage {
			isReplaced = true
		}
	}
	test.Assert(t, isReplaced, "blockchain is kept whereas it should be replaced")
}

func Test_Update_NeighborNewBlockTimestampIsInvalid_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, 0, registry, 1, synchronizer, logger)

	type args struct {
		firstBlockTimestamp  int64
		secondBlockTimestamp int64
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "SecondTimestampBeforeTheFirstOne",
			args: args{
				firstBlockTimestamp:  1,
				secondBlockTimestamp: 0,
			},
		},
		{
			name: "BlockMissing",
			args: args{
				firstBlockTimestamp:  0,
				secondBlockTimestamp: 2,
			},
		},
		{
			name: "SameZeroedTimestamp",
			args: args{
				firstBlockTimestamp:  0,
				secondBlockTimestamp: 0,
			},
		},
		{
			name: "SameNonZeroTimestamp",
			args: args{
				firstBlockTimestamp:  1,
				secondBlockTimestamp: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
				blockResponse1 := protocoltest.NewRewardedBlockResponse([32]byte{}, tt.args.firstBlockTimestamp)
				block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
				hash, _ := block1.Hash()
				blockResponse2 := protocoltest.NewRewardedBlockResponse(hash, tt.args.secondBlockTimestamp)
				return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
			}

			// Act
			blockchain.Update(1)

			// Assert
			var isKept bool
			var isExplicitMessageLogged bool
			for _, call := range logger.DebugCalls() {
				expectedMessage := "neighbor block timestamp is invalid"
				if call.Msg == blockchainKeptMessage {
					isKept = true
				} else if strings.Contains(call.Msg, expectedMessage) {
					isExplicitMessageLogged = true
				}
			}
			test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
			test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
		})
	}
}

func Test_Update_NeighborNewBlockTimestampIsInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewRewardedBlockResponse([32]byte{}, 1)
		block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
		hash, _ := block1.Hash()
		blockResponse2 := protocoltest.NewRewardedBlockResponse(hash, 2)
		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, 0, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(1)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := "neighbor block timestamp is in the future"
		if call.Msg == blockchainKeptMessage {
			isKept = true
		} else if strings.Contains(call.Msg, expectedMessage) {
			isExplicitMessageLogged = true
		}
	}
	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
}

func Test_Update_NeighborNewBlockTransactionFeeIsTooLow_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	wallet, _ := encryption.NewWallet("", "", "", test.PrivateKey)
	address := wallet.Address()
	var invalidTransactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	transactionRequest := protocoltest.NewSignedTransactionRequest(invalidTransactionFee, "A", address, privateKey, publicKey, 3, 1)
	transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
	transactionResponse := transaction.GetResponse()
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
		hash, _ := block1.Hash()
		var block2Timestamp int64 = 1
		transactions := []*network.TransactionResponse{
			transactionResponse,
			validation.NewRewardTransaction(address, block2Timestamp, 0),
		}
		blockResponse2 := verification.NewBlockResponse(block2Timestamp, hash, transactions, []string{address}, nil)
		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	var minimalTransactionFee uint64 = 1
	blockchain := verification.NewBlockchain(genesisTransaction, minimalTransactionFee, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(1)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := fmt.Sprintf("a neighbor block transaction fee is too low, fee: %d, minimal fee: %d", invalidTransactionFee, minimalTransactionFee)
		if call.Msg == blockchainKeptMessage {
			isKept = true
		} else if strings.Contains(call.Msg, expectedMessage) {
			isExplicitMessageLogged = true
		}
	}
	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
}

func Test_Update_NeighborNewBlockTransactionTimestampIsTooFarInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	wallet, _ := encryption.NewWallet("", "", "", test.PrivateKey)
	address := wallet.Address()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	transactionRequest := protocoltest.NewSignedTransactionRequest(transactionFee, "A", address, privateKey, publicKey, 3, 1)
	transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
	transactionResponse := transaction.GetResponse()
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
		hash, _ := block1.Hash()
		var block2Timestamp int64 = 1
		transactions := []*network.TransactionResponse{
			transactionResponse,
			validation.NewRewardTransaction(address, block2Timestamp, 0),
		}
		blockResponse2 := verification.NewBlockResponse(block2Timestamp, hash, transactions, []string{address}, nil)
		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, transactionFee, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(1)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := fmt.Sprintf("a neighbor block transaction timestamp is too far in the future, transaction: %v", transactionResponse)
		if call.Msg == blockchainKeptMessage {
			isKept = true
		} else if strings.Contains(call.Msg, expectedMessage) {
			isExplicitMessageLogged = true
		}
	}
	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
}

func Test_Update_NeighborNewBlockTransactionTimestampIsTooOld_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	wallet, _ := encryption.NewWallet("", "", "", test.PrivateKey)
	address := wallet.Address()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	transactionRequest := protocoltest.NewSignedTransactionRequest(transactionFee, "A", address, privateKey, publicKey, 0, 1)
	transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
	transactionResponse := transaction.GetResponse()
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
		hash1, _ := block1.Hash()
		blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, 1)
		block2, _ := verification.NewBlockFromResponse(blockResponse2, nil)
		hash2, _ := block2.Hash()
		var block3Timestamp int64 = 2
		transactions := []*network.TransactionResponse{
			transactionResponse,
			validation.NewRewardTransaction(address, block3Timestamp, 0),
		}
		blockResponse3 := verification.NewBlockResponse(block3Timestamp, hash2, transactions, []string{address}, nil)
		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, transactionFee, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(2)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := fmt.Sprintf("a neighbor block transaction timestamp is too old, transaction: %v", transactionResponse)
		if call.Msg == blockchainKeptMessage {
			isKept = true
		} else if strings.Contains(call.Msg, expectedMessage) {
			isExplicitMessageLogged = true
		}
	}
	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
}
