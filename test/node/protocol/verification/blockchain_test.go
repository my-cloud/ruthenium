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
	"math"
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
	blockchain := verification.NewBlockchain(0, nil, 1, 1, 1, 0, registry, 1, synchronizer, logger)

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
	blockchain := verification.NewBlockchain(0, nil, 1, 1, 1, 0, registry, 1, synchronizer, logger)

	// Act
	blocks := blockchain.Blocks()

	// Assert
	test.Assert(t, len(blocks) == 1, "blocks don't contain a single block")
}

func Test_UtxosByAddress_UnknownAddress_ReturnsNil(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	genesisValidatorAddress := ""
	genesisTransaction, _ := validation.NewGenesisTransaction(genesisValidatorAddress, 0, 0)
	blockchain := verification.NewBlockchain(0, genesisTransaction, 1, 1, 1, 0, nil, 1, nil, logger)

	// Act
	utxos := blockchain.UtxosByAddress(genesisValidatorAddress)

	// Assert
	test.Assert(t, utxos == nil, "utxos list is not nil whereas it should be")
}

func Test_UtxosByAddress_GenesisValidator_ReturnsGenesisUtxo(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	var genesisTransactionValue uint64 = 10
	genesisValidatorAddress := ""
	genesisTransaction, _ := validation.NewGenesisTransaction(genesisValidatorAddress, 0, genesisTransactionValue)
	blockchain := verification.NewBlockchain(0, genesisTransaction, 1, 1, 1, 0, registry, 1, synchronizer, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	utxos := blockchain.UtxosByAddress(genesisValidatorAddress)

	// Assert
	test.Assert(t, utxos[0].Value == genesisTransactionValue, "utxo amount is not the genesis amount whereas it should be")
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
	address := test.Address
	genesisTransaction, _ := validation.NewGenesisTransaction(address, 0, 1)
	blockchain := verification.NewBlockchain(0, genesisTransaction, 1, 1, 1, 0, registry, 1, synchronizer, logger)
	_ = blockchain.AddBlock(1, nil, nil)
	_ = blockchain.AddBlock(2, nil, nil)
	blocks := blockchain.LastBlocks(0)
	genesisBlockHash := blocks[1].PreviousHash
	blockResponse1 := protocoltest.NewRewardedBlockResponse(genesisBlockHash, 1)
	block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
	hash1, _ := block1.Hash()
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var now int64 = 2
	var genesisOutputIndex uint16 = 0
	transactionRequest := protocoltest.NewSignedTransactionRequest(0, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, 3, 0)
	transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
	transactionResponse := transaction.GetResponse()
	rewardTransaction, _ := validation.NewRewardTransaction(address, now, 0)
	transactions := []*network.TransactionResponse{
		transactionResponse,
		rewardTransaction,
	}
	blockResponse2 := verification.NewBlockResponse(2, hash1, transactions, []string{address}, nil)
	block2, _ := verification.NewBlockFromResponse(blockResponse2, nil)
	hash2, _ := block2.Hash()
	blockResponse3 := protocoltest.NewRewardedBlockResponse(hash2, 3)
	block3, _ := verification.NewBlockFromResponse(blockResponse3, nil)
	hash3, _ := block3.Hash()
	blockResponse4 := protocoltest.NewRewardedBlockResponse(hash3, 4)
	lastBlocksResponses := []*network.BlockResponse{blockResponse2, blockResponse3, blockResponse4}
	neighborMock.GetLastBlocksFunc = func(uint64) ([]*network.BlockResponse, error) { return lastBlocksResponses, nil }
	blockResponses := []*network.BlockResponse{blocks[0], blockResponse1, blockResponse2, blockResponse3, blockResponse4}
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) { return blockResponses, nil }

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
	blockchain := verification.NewBlockchain(0, nil, 1, 1, 1, 0, registry, 1, synchronizer, logger)

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
	blockchain := verification.NewBlockchain(0, nil, 1, 1, 1, 0, registry, 1, synchronizer, logger)

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

func Test_Update_NeighborNewBlockTransactionFeeIsNegative_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	address := test.Address
	var invalidTransactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var now int64 = 2
	var incomeLimit uint64 = 1
	genesisValue := 2 * incomeLimit
	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisValue)
	block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
	hash1, _ := block1.Hash()
	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-1)
	block2, _ := verification.NewBlockFromResponse(blockResponse2, nil)
	hash2, _ := block2.Hash()
	genesisTransaction := blockResponse1.Transactions[0]
	var genesisOutputIndex uint16 = 0
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, invalidTransactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, 3, genesisValue)
	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	invalidTransactionResponse := invalidTransaction.GetResponse()
	rewardTransaction, _ := validation.NewRewardTransaction(address, now, 1)
	transactions := []*network.TransactionResponse{
		invalidTransactionResponse,
		rewardTransaction,
	}
	blockResponse3 := verification.NewBlockResponse(now, hash2, transactions, []string{address}, nil)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	var minimalTransactionFee uint64 = 1
	blockchain := verification.NewBlockchain(0, nil, 1, 1, incomeLimit, minimalTransactionFee, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(now)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := "transaction fee is negative"
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
	address := test.Address
	var invalidTransactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var now int64 = 2
	var genesisValue uint64 = 1
	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisValue)
	block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
	hash1, _ := block1.Hash()
	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-1)
	block2, _ := verification.NewBlockFromResponse(blockResponse2, nil)
	hash2, _ := block2.Hash()
	genesisTransaction := blockResponse1.Transactions[0]
	var genesisOutputIndex uint16 = 0
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisValue, invalidTransactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, 3, genesisValue)
	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	invalidTransactionResponse := invalidTransaction.GetResponse()
	rewardTransaction, _ := validation.NewRewardTransaction(address, now, 1)
	transactions := []*network.TransactionResponse{
		invalidTransactionResponse,
		rewardTransaction,
	}
	blockResponse3 := verification.NewBlockResponse(now, hash2, transactions, []string{address}, nil)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	var minimalTransactionFee uint64 = 1
	blockchain := verification.NewBlockchain(0, nil, 1, 1, genesisValue, minimalTransactionFee, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(now)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := "transaction fee is too low"
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
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var now int64 = 2
	genesisValue := uint64(math.Pow(2, float64(now)))
	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisValue)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := blockResponse1.Transactions[0]
	transactionValue := uint64(math.Pow(float64(genesisValue), -float64(now)))
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(transactionValue, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, 4, transactionValue)
	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	invalidTransactionResponse := invalidTransaction.GetResponse()
	block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
	hash1, _ := block1.Hash()
	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-1)
	block2, _ := verification.NewBlockFromResponse(blockResponse2, nil)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := validation.NewRewardTransaction(address, now, 0)
	transactions := []*network.TransactionResponse{
		invalidTransactionResponse,
		rewardTransaction,
	}
	blockResponse3 := verification.NewBlockResponse(now, hash2, transactions, []string{address}, nil)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(0, nil, 1, 1, 1, transactionFee, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(now)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := fmt.Sprintf("a neighbor block transaction timestamp is too far in the future, transaction: %v", invalidTransactionResponse)
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
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var now int64 = 2
	genesisValue := uint64(math.Pow(2, float64(now)))
	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisValue)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := blockResponse1.Transactions[0]
	transactionValue := uint64(math.Pow(float64(genesisValue), -float64(now)))
	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(transactionValue, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, 0, transactionValue)
	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
	invalidTransactionResponse := invalidTransaction.GetResponse()
	block1, _ := verification.NewBlockFromResponse(blockResponse1, nil)
	hash1, _ := block1.Hash()
	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-1)
	block2, _ := verification.NewBlockFromResponse(blockResponse2, nil)
	hash2, _ := block2.Hash()
	var block3Timestamp int64 = 2
	rewardTransaction, _ := validation.NewRewardTransaction(address, block3Timestamp, 0)
	transactions := []*network.TransactionResponse{
		invalidTransactionResponse,
		rewardTransaction,
	}
	blockResponse3 := verification.NewBlockResponse(block3Timestamp, hash2, transactions, []string{address}, nil)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(networktest.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(0, nil, 1, 1, 1, transactionFee, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(now)

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := fmt.Sprintf("a neighbor block transaction timestamp is too old, transaction: %v", invalidTransactionResponse)
		if call.Msg == blockchainKeptMessage {
			isKept = true
		} else if strings.Contains(call.Msg, expectedMessage) {
			isExplicitMessageLogged = true
		}
	}
	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
}
