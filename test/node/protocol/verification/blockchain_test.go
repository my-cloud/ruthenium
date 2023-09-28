package verification

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/config"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocoltest"
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
	settings := config.Settings{
		GenesisAmountInParticles:    1,
		HalfLifeInDays:              1,
		IncomeBaseInParticles:       1,
		IncomeLimitInParticles:      1,
		MinimalTransactionFee:       0,
		ValidationIntervalInSeconds: 1,
	}
	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)

	// Act
	err := blockchain.AddBlock(0, nil, nil)

	// Assert
	test.Assert(t, err == nil, "error is returned whereas it should not")
}

func Test_Blocks_BlocksCountLimitSetToZero_ReturnsNil(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	settings := config.Settings{
		GenesisAmountInParticles:    1,
		HalfLifeInDays:              1,
		IncomeBaseInParticles:       1,
		IncomeLimitInParticles:      1,
		MinimalTransactionFee:       0,
		ValidationIntervalInSeconds: 1,
	}
	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)

	// Act
	blocks := blockchain.Blocks(0)

	// Assert
	test.Assert(t, len(blocks) == 0, "blocks should be empty")
}

func Test_Blocks_BlocksCountLimitSetToOne_ReturnsOneBlock(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	var expectedBlocksCount uint64 = 1
	var validationInterval int64 = 1
	settings := config.Settings{
		BlocksCountLimit:            expectedBlocksCount,
		GenesisAmountInParticles:    1,
		HalfLifeInDays:              1,
		IncomeBaseInParticles:       1,
		IncomeLimitInParticles:      1,
		MinimalTransactionFee:       0,
		ValidationIntervalInSeconds: validationInterval,
	}
	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, nil, nil)

	// Act
	blocksBytes := blockchain.Blocks(0)

	// Assert
	var blocks []*verification.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	actualBlocksCount := uint64(len(blocks))
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_Blocks_BlocksCountLimitSetToTwo_ReturnsTwoBlocks(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	var expectedBlocksCount uint64 = 2
	var validationInterval int64 = 1
	settings := config.Settings{
		BlocksCountLimit:            expectedBlocksCount,
		GenesisAmountInParticles:    1,
		HalfLifeInDays:              1,
		IncomeBaseInParticles:       1,
		IncomeLimitInParticles:      1,
		MinimalTransactionFee:       0,
		ValidationIntervalInSeconds: validationInterval,
	}
	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, nil, nil)

	// Act
	blocksBytes := blockchain.Blocks(0)

	// Assert
	var blocks []*verification.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	actualBlocksCount := uint64(len(blocks))
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_Blocks_StartingBlockHeightGreaterThanBlocksLength_ReturnsNil(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	var validationInterval int64 = 1
	settings := config.Settings{
		BlocksCountLimit:            1,
		GenesisAmountInParticles:    1,
		HalfLifeInDays:              1,
		IncomeBaseInParticles:       1,
		IncomeLimitInParticles:      1,
		MinimalTransactionFee:       0,
		ValidationIntervalInSeconds: validationInterval,
	}
	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)

	// Act
	blocks := blockchain.Blocks(1)

	// Assert
	expectedBlocksCount := 0
	actualBlocksCount := len(blocks)
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_UtxosByAddress_UnknownAddress_ReturnsNil(t *testing.T) {
	// Arrange
	logger := logtest.NewLoggerMock()
	genesisValidatorAddress := ""
	var genesisAmount uint64 = 0
	settings := config.Settings{
		GenesisAmountInParticles:    genesisAmount,
		HalfLifeInDays:              1,
		IncomeBaseInParticles:       1,
		IncomeLimitInParticles:      1,
		MinimalTransactionFee:       0,
		ValidationIntervalInSeconds: 1,
	}
	blockchain := verification.NewBlockchain(nil, settings, nil, logger)

	// Act
	utxos := blockchain.UtxosByAddress(genesisValidatorAddress)

	// Assert
	test.Assert(t, len(utxos) == 0, "utxos should be empty")
}

func Test_UtxosByAddress_UtxoExists_ReturnsUtxo(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	var validationInterval int64 = 1
	settings := config.Settings{
		GenesisAmountInParticles:    1,
		HalfLifeInDays:              1,
		IncomeBaseInParticles:       1,
		IncomeLimitInParticles:      1,
		MinimalTransactionFee:       0,
		ValidationIntervalInSeconds: validationInterval,
	}
	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
	registeredAddress := ""
	var expectedValue uint64 = 1
	var genesisTimestamp int64 = 0
	transaction, _ := validation.NewRewardTransaction(registeredAddress, true, genesisTimestamp+validationInterval, expectedValue)
	transactions := []*validation.Transaction{transaction}
	transactionsBytes, _ := json.Marshal(transactions)
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, transactionsBytes, []string{registeredAddress})
	_ = blockchain.AddBlock(genesisTimestamp+2*validationInterval, nil, nil)

	// Act
	utxos := blockchain.UtxosByAddress(registeredAddress)

	// Assert
	actualValue := utxos[0].Value
	test.Assert(t, actualValue == expectedValue, fmt.Sprintf("utxo amount is %d whereas it should be %d", actualValue, expectedValue))
}

//
//func Test_Update_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
//	// Arrange
//	registry := new(protocoltest.RegistryMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := logtest.NewLoggerMock()
//	neighborMock := new(networktest.NeighborMock)
//	neighborMock.TargetFunc = func() string {
//		return "neighbor"
//	}
//	synchronizer := new(networktest.SynchronizerMock)
//	synchronizer.NeighborsFunc = func() []network.Neighbor {
//		return []network.Neighbor{neighborMock}
//	}
//	var validationIntervalInSeconds int64 = 1
//	settings := config.Settings{
//		GenesisAmountInParticles:    1,
//		HalfLifeInDays:              1,
//		IncomeBaseInParticles:       1,
//		IncomeLimitInParticles:      1,
//		MinimalTransactionFee:       0,
//		ValidationIntervalInSeconds: 1,
//	}
//	validationTimestamp := validationIntervalInSeconds * time.Second.Nanoseconds()
//	now := 5 * validationTimestamp
//	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
//	_ = blockchain.AddBlock(now-5*validationTimestamp, nil, nil)
//	_ = blockchain.AddBlock(now-4*validationTimestamp, nil, nil)
//	blocks := blockchain.AllBlocks()
//	genesisBlockHash := blocks[1].PreviousHash
//	blockResponse1 := protocoltest.NewRewardedBlockResponse(genesisBlockHash, now-4*validationTimestamp)
//	block1, _ := verification.NewBlockFromResponse(blockResponse1)
//	hash1, _ := block1.Hash()
//	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-3*validationTimestamp)
//	block2, _ := verification.NewBlockFromResponse(blockResponse2)
//	hash2, _ := block2.Hash()
//	blockResponse3 := protocoltest.NewRewardedBlockResponse(hash2, now-2*validationTimestamp)
//	block3, _ := verification.NewBlockFromResponse(blockResponse3)
//	hash3, _ := block3.Hash()
//	blockResponse4 := protocoltest.NewRewardedBlockResponse(hash3, now-validationTimestamp)
//	blockResponses := []*network.BlockResponse{blocks[0], blockResponse1, blockResponse2, blockResponse3, blockResponse4}
//	neighborMock.GetBlocksFunc = func(uint64) ([]*network.BlockResponse, error) { return blockResponses, nil }
//
//	// Act
//	blockchain.Update(now)
//
//	// Assert
//	var isReplaced bool
//	for _, call := range logger.DebugCalls() {
//		if call.Msg == blockchainReplacedMessage {
//			isReplaced = true
//		}
//	}
//	test.Assert(t, isReplaced, "blockchain is kept whereas it should be replaced")
//}
//
//func Test_Update_NeighborNewBlockTimestampIsInvalid_IsNotReplaced(t *testing.T) {
//	// Arrange
//	registry := new(protocoltest.RegistryMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := logtest.NewLoggerMock()
//	neighborMock := new(networktest.NeighborMock)
//	neighborMock.TargetFunc = func() string {
//		return "neighbor"
//	}
//	synchronizer := new(networktest.SynchronizerMock)
//	synchronizer.NeighborsFunc = func() []network.Neighbor {
//		return []network.Neighbor{neighborMock}
//	}
//	settings := config.Settings{
//		GenesisAmountInParticles:    1,
//		HalfLifeInDays:              1,
//		IncomeBaseInParticles:       1,
//		IncomeLimitInParticles:      1,
//		MinimalTransactionFee:       0,
//		ValidationIntervalInSeconds: 1,
//	}
//	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
//	_ = blockchain.AddBlock(0, nil, nil)
//
//	type args struct {
//		firstBlockTimestamp  int64
//		secondBlockTimestamp int64
//	}
//	tests := []struct {
//		name string
//		args args
//		want []int
//	}{
//		{
//			name: "SecondTimestampBeforeTheFirstOne",
//			args: args{
//				firstBlockTimestamp:  1,
//				secondBlockTimestamp: 0,
//			},
//		},
//		{
//			name: "BlockMissing",
//			args: args{
//				firstBlockTimestamp:  0,
//				secondBlockTimestamp: 2,
//			},
//		},
//		{
//			name: "SameZeroedTimestamp",
//			args: args{
//				firstBlockTimestamp:  0,
//				secondBlockTimestamp: 0,
//			},
//		},
//		{
//			name: "SameNonZeroTimestamp",
//			args: args{
//				firstBlockTimestamp:  1,
//				secondBlockTimestamp: 1,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			neighborMock.GetBlocksFunc = func(uint64) ([]*network.BlockResponse, error) {
//				blockResponse1 := protocoltest.NewRewardedBlockResponse([32]byte{}, tt.args.firstBlockTimestamp)
//				block1, _ := verification.NewBlockFromResponse(blockResponse1)
//				hash, _ := block1.Hash()
//				blockResponse2 := protocoltest.NewRewardedBlockResponse(hash, tt.args.secondBlockTimestamp)
//				return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
//			}
//
//			// Act
//			blockchain.Update(1)
//
//			// Assert
//			var isKept bool
//			var isExplicitMessageLogged bool
//			for _, call := range logger.DebugCalls() {
//				expectedMessage := "neighbor block timestamp is invalid"
//				if call.Msg == blockchainKeptMessage {
//					isKept = true
//				} else if strings.Contains(call.Msg, expectedMessage) {
//					isExplicitMessageLogged = true
//				}
//			}
//			test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
//			test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
//		})
//	}
//}
//
//func Test_Update_NeighborNewBlockTimestampIsInTheFuture_IsNotReplaced(t *testing.T) {
//	// Arrange
//	registry := new(protocoltest.RegistryMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := logtest.NewLoggerMock()
//	neighborMock := new(networktest.NeighborMock)
//	var validationIntervalInSeconds int64 = 1
//	validationTimestamp := validationIntervalInSeconds * time.Second.Nanoseconds()
//	now := validationTimestamp
//	neighborMock.GetBlocksFunc = func(uint64) ([]*network.BlockResponse, error) {
//		blockResponse1 := protocoltest.NewRewardedBlockResponse([32]byte{}, now)
//		block1, _ := verification.NewBlockFromResponse(blockResponse1)
//		hash, _ := block1.Hash()
//		blockResponse2 := protocoltest.NewRewardedBlockResponse(hash, now+validationTimestamp)
//		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
//	}
//	neighborMock.TargetFunc = func() string {
//		return "neighbor"
//	}
//	synchronizer := new(networktest.SynchronizerMock)
//	synchronizer.NeighborsFunc = func() []network.Neighbor {
//		return []network.Neighbor{neighborMock}
//	}
//	settings := config.Settings{
//		GenesisAmountInParticles:    1,
//		HalfLifeInDays:              1,
//		IncomeBaseInParticles:       1,
//		IncomeLimitInParticles:      1,
//		MinimalTransactionFee:       0,
//		ValidationIntervalInSeconds: validationIntervalInSeconds,
//	}
//	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
//	_ = blockchain.AddBlock(0, nil, nil)
//
//	// Act
//	blockchain.Update(now)
//
//	// Assert
//	var isKept bool
//	var isExplicitMessageLogged bool
//	for _, call := range logger.DebugCalls() {
//		expectedMessage := "neighbor block timestamp is in the future"
//		if call.Msg == blockchainKeptMessage {
//			isKept = true
//		} else if strings.Contains(call.Msg, expectedMessage) {
//			isExplicitMessageLogged = true
//		}
//	}
//	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
//	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
//}
//
//func Test_Update_NeighborNewBlockTransactionFeeIsNegative_IsNotReplaced(t *testing.T) {
//	// Arrange
//	registry := new(protocoltest.RegistryMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := logtest.NewLoggerMock()
//	neighborMock := new(networktest.NeighborMock)
//	address := test.Address
//	var invalidTransactionFee uint64 = 0
//	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
//	publicKey := encryption.NewPublicKey(privateKey)
//	var validationIntervalInSeconds int64 = 1
//	validationTimestamp := validationIntervalInSeconds * time.Second.Nanoseconds()
//	now := 2 * validationTimestamp
//	var incomeLimit uint64 = 1
//	genesisAmount := 2 * incomeLimit
//	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisAmount)
//	block1, _ := verification.NewBlockFromResponse(blockResponse1)
//	hash1, _ := block1.Hash()
//	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-validationTimestamp)
//	block2, _ := verification.NewBlockFromResponse(blockResponse2)
//	hash2, _ := block2.Hash()
//	genesisTransaction := blockResponse1.Transactions[0]
//	var genesisOutputIndex uint16 = 0
//	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisAmount, invalidTransactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now+validationTimestamp, genesisAmount)
//	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
//	invalidTransactionResponse := invalidTransaction.GetResponse()
//	rewardTransaction, _ := validation.NewRewardTransaction(address, false, now, 1)
//	transactions := []*network.TransactionResponse{
//		invalidTransactionResponse,
//		rewardTransaction,
//	}
//	blockResponse3 := verification.NewBlockResponse(now, hash2, transactions, []string{address}, nil)
//	neighborMock.GetBlocksFunc = func(uint64) ([]*network.BlockResponse, error) {
//		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
//	}
//	neighborMock.TargetFunc = func() string {
//		return "neighbor"
//	}
//	synchronizer := new(networktest.SynchronizerMock)
//	synchronizer.NeighborsFunc = func() []network.Neighbor {
//		return []network.Neighbor{neighborMock}
//	}
//	settings := config.Settings{
//		GenesisAmountInParticles:         genesisAmount,
//		HalfLifeInDays:                   1,
//		IncomeBaseInParticles:            1,
//		IncomeLimitInParticles:           incomeLimit,
//		MaxOutboundsCount:                1,
//		MinimalTransactionFee:            1,
//		SynchronizationIntervalInSeconds: 1,
//		ValidationIntervalInSeconds:      1,
//		VerificationsCountPerValidation:  1,
//	}
//	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
//	_ = blockchain.AddBlock(0, nil, nil)
//
//	// Act
//	blockchain.Update(now)
//
//	// Assert
//	var isKept bool
//	var isExplicitMessageLogged bool
//	for _, call := range logger.DebugCalls() {
//		expectedMessage := "transaction fee is negative"
//		if call.Msg == blockchainKeptMessage {
//			isKept = true
//		} else if strings.Contains(call.Msg, expectedMessage) {
//			isExplicitMessageLogged = true
//		}
//	}
//	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
//	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
//}
//
//func Test_Update_NeighborNewBlockTransactionFeeIsTooLow_IsNotReplaced(t *testing.T) {
//	// Arrange
//	registry := new(protocoltest.RegistryMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := logtest.NewLoggerMock()
//	neighborMock := new(networktest.NeighborMock)
//	address := test.Address
//	var invalidTransactionFee uint64 = 0
//	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
//	publicKey := encryption.NewPublicKey(privateKey)
//	var validationIntervalInSeconds int64 = 1
//	validationTimestamp := validationIntervalInSeconds * time.Second.Nanoseconds()
//	now := 2 * validationTimestamp
//	var genesisAmount uint64 = 1
//	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisAmount)
//	block1, _ := verification.NewBlockFromResponse(blockResponse1)
//	hash1, _ := block1.Hash()
//	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-validationTimestamp)
//	block2, _ := verification.NewBlockFromResponse(blockResponse2)
//	hash2, _ := block2.Hash()
//	genesisTransaction := blockResponse1.Transactions[0]
//	var genesisOutputIndex uint16 = 0
//	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisAmount, invalidTransactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now+validationTimestamp, genesisAmount)
//	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
//	invalidTransactionResponse := invalidTransaction.GetResponse()
//	rewardTransaction, _ := validation.NewRewardTransaction(address, false, now, 1)
//	transactions := []*network.TransactionResponse{
//		invalidTransactionResponse,
//		rewardTransaction,
//	}
//	blockResponse3 := verification.NewBlockResponse(now, hash2, transactions, []string{address}, nil)
//	neighborMock.GetBlocksFunc = func(uint64) ([]*network.BlockResponse, error) {
//		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
//	}
//	neighborMock.TargetFunc = func() string {
//		return "neighbor"
//	}
//	synchronizer := new(networktest.SynchronizerMock)
//	synchronizer.NeighborsFunc = func() []network.Neighbor {
//		return []network.Neighbor{neighborMock}
//	}
//	settings := config.Settings{
//		GenesisAmountInParticles:    genesisAmount,
//		HalfLifeInDays:              1,
//		IncomeBaseInParticles:       1,
//		IncomeLimitInParticles:      1,
//		MinimalTransactionFee:       1,
//		ValidationIntervalInSeconds: 1,
//	}
//	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
//	_ = blockchain.AddBlock(0, nil, nil)
//
//	// Act
//	blockchain.Update(now)
//
//	// Assert
//	var isKept bool
//	var isExplicitMessageLogged bool
//	for _, call := range logger.DebugCalls() {
//		expectedMessage := "transaction fee is too low"
//		if call.Msg == blockchainKeptMessage {
//			isKept = true
//		} else if strings.Contains(call.Msg, expectedMessage) {
//			isExplicitMessageLogged = true
//		}
//	}
//	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
//	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
//}
//
//func Test_Update_NeighborNewBlockTransactionTimestampIsTooFarInTheFuture_IsNotReplaced(t *testing.T) {
//	// Arrange
//	registry := new(protocoltest.RegistryMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := logtest.NewLoggerMock()
//	neighborMock := new(networktest.NeighborMock)
//	address := test.Address
//	var transactionFee uint64 = 0
//	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
//	publicKey := encryption.NewPublicKey(privateKey)
//	var validationIntervalInSeconds int64 = 1
//	validationTimestamp := validationIntervalInSeconds * time.Second.Nanoseconds()
//	now := 2 * validationTimestamp
//	var genesisAmount uint64 = 1
//	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisAmount)
//	var genesisOutputIndex uint16 = 0
//	genesisTransaction := blockResponse1.Transactions[0]
//	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisAmount, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, now+2*validationTimestamp, genesisAmount)
//	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
//	invalidTransactionResponse := invalidTransaction.GetResponse()
//	block1, _ := verification.NewBlockFromResponse(blockResponse1)
//	hash1, _ := block1.Hash()
//	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-validationTimestamp)
//	block2, _ := verification.NewBlockFromResponse(blockResponse2)
//	hash2, _ := block2.Hash()
//	rewardTransaction, _ := validation.NewRewardTransaction(address, false, now, 0)
//	transactions := []*network.TransactionResponse{
//		invalidTransactionResponse,
//		rewardTransaction,
//	}
//	blockResponse3 := verification.NewBlockResponse(now, hash2, transactions, []string{address}, nil)
//	neighborMock.GetBlocksFunc = func(uint64) ([]*network.BlockResponse, error) {
//		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
//	}
//	neighborMock.TargetFunc = func() string {
//		return "neighbor"
//	}
//	synchronizer := new(networktest.SynchronizerMock)
//	synchronizer.NeighborsFunc = func() []network.Neighbor {
//		return []network.Neighbor{neighborMock}
//	}
//	settings := config.Settings{
//		GenesisAmountInParticles:    genesisAmount,
//		HalfLifeInDays:              1,
//		IncomeBaseInParticles:       1,
//		IncomeLimitInParticles:      1,
//		MinimalTransactionFee:       transactionFee,
//		ValidationIntervalInSeconds: 1,
//	}
//	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
//	_ = blockchain.AddBlock(0, nil, nil)
//
//	// Act
//	blockchain.Update(now)
//
//	// Assert
//	var isKept bool
//	var isExplicitMessageLogged bool
//	for _, call := range logger.DebugCalls() {
//		expectedMessage := fmt.Sprintf("a neighbor block transaction timestamp is too far in the future, transaction: %v", invalidTransactionResponse)
//		if call.Msg == blockchainKeptMessage {
//			isKept = true
//		} else if strings.Contains(call.Msg, expectedMessage) {
//			isExplicitMessageLogged = true
//		}
//	}
//	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
//	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
//}
//
//func Test_Update_NeighborNewBlockTransactionTimestampIsTooOld_IsNotReplaced(t *testing.T) {
//	// Arrange
//	registry := new(protocoltest.RegistryMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := logtest.NewLoggerMock()
//	neighborMock := new(networktest.NeighborMock)
//	address := test.Address
//	var transactionFee uint64 = 0
//	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
//	publicKey := encryption.NewPublicKey(privateKey)
//	var validationIntervalInSeconds int64 = 1
//	validationTimestamp := validationIntervalInSeconds * time.Second.Nanoseconds()
//	now := 2 * validationTimestamp
//	var genesisAmount uint64 = 1
//	blockResponse1 := protocoltest.NewGenesisBlockResponse(address, genesisAmount)
//	var genesisOutputIndex uint16 = 0
//	genesisTransaction := blockResponse1.Transactions[0]
//	invalidTransactionRequest := protocoltest.NewSignedTransactionRequest(genesisAmount, transactionFee, "A", genesisTransaction, genesisOutputIndex, privateKey, publicKey, 0, genesisAmount)
//	invalidTransaction, _ := validation.NewTransactionFromRequest(&invalidTransactionRequest)
//	invalidTransactionResponse := invalidTransaction.GetResponse()
//	block1, _ := verification.NewBlockFromResponse(blockResponse1)
//	hash1, _ := block1.Hash()
//	blockResponse2 := protocoltest.NewRewardedBlockResponse(hash1, now-validationTimestamp)
//	block2, _ := verification.NewBlockFromResponse(blockResponse2)
//	hash2, _ := block2.Hash()
//	rewardTransaction, _ := validation.NewRewardTransaction(address, false, now, 0)
//	transactions := []*network.TransactionResponse{
//		invalidTransactionResponse,
//		rewardTransaction,
//	}
//	blockResponse3 := verification.NewBlockResponse(now, hash2, transactions, []string{address}, nil)
//	neighborMock.GetBlocksFunc = func(uint64) ([]*network.BlockResponse, error) {
//		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
//	}
//	neighborMock.TargetFunc = func() string {
//		return "neighbor"
//	}
//	synchronizer := new(networktest.SynchronizerMock)
//	synchronizer.NeighborsFunc = func() []network.Neighbor {
//		return []network.Neighbor{neighborMock}
//	}
//	settings := config.Settings{
//		GenesisAmountInParticles:    genesisAmount,
//		HalfLifeInDays:              1,
//		IncomeBaseInParticles:       1,
//		IncomeLimitInParticles:      1,
//		MinimalTransactionFee:       transactionFee,
//		ValidationIntervalInSeconds: 1,
//	}
//	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
//	_ = blockchain.AddBlock(0, nil, nil)
//
//	// Act
//	blockchain.Update(now)
//
//	// Assert
//	var isKept bool
//	var isExplicitMessageLogged bool
//	for _, call := range logger.DebugCalls() {
//		expectedMessage := fmt.Sprintf("a neighbor block transaction timestamp is too old, transaction: %v", invalidTransactionResponse)
//		if call.Msg == blockchainKeptMessage {
//			isKept = true
//		} else if strings.Contains(call.Msg, expectedMessage) {
//			isExplicitMessageLogged = true
//		}
//	}
//	test.Assert(t, isKept, "blockchain is replaced whereas it should be kept")
//	test.Assert(t, isExplicitMessageLogged, "no explicit message is logged whereas it should be")
//}
