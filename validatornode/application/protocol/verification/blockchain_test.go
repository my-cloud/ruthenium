package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol"
	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"github.com/my-cloud/ruthenium/validatornode/presentation"
)

const (
	blockchainReplacedMessage = "verification done: blockchain replaced"
	blockchainKeptMessage     = "verification done: blockchain kept"
)

func Test_AddBlock_ValidParameters_NoErrorReturned(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	settings := new(protocol.SettingsProviderMock)
	utxosManagerMock := new(protocol.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)

	// Act
	err := blockchain.AddBlock(0, nil, nil)

	// Assert
	test.Assert(t, err == nil, "error is returned whereas it should not")
}

func Test_Blocks_BlocksCountLimitSetToZero_ReturnsEmptyArray(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	settings := new(protocol.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 0 }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)

	// Act
	blocksBytes := blockchain.Blocks(0)

	// Assert
	var blocks []*ledger.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	test.Assert(t, len(blocks) == 0, "blocks should be empty")
}

func Test_Blocks_BlocksCountLimitSetToOne_ReturnsOneBlock(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	var expectedBlocksCount uint64 = 1
	settings := new(protocol.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return expectedBlocksCount }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	var validationInterval int64 = 1
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, nil, nil)

	// Act
	blocksBytes := blockchain.Blocks(0)

	// Assert
	var blocks []*ledger.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	actualBlocksCount := uint64(len(blocks))
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_Blocks_BlocksCountLimitSetToTwo_ReturnsTwoBlocks(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	var expectedBlocksCount uint64 = 2
	settings := new(protocol.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return expectedBlocksCount }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	var validationInterval int64 = 1
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, nil, nil)

	// Act
	blocksBytes := blockchain.Blocks(0)

	// Assert
	var blocks []*ledger.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	actualBlocksCount := uint64(len(blocks))
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_Blocks_StartingBlockHeightGreaterThanBlocksLength_ReturnsEmptyArray(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	var blocksCount uint64 = 1
	settings := new(protocol.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return blocksCount }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)

	// Act
	blocksBytes := blockchain.Blocks(1)

	// Assert
	expectedBlocksCount := 0
	var blocks []*ledger.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	actualBlocksCount := len(blocks)
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_FirstBlockTimestamp_BlockchainIsEmpty_Returns0(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	settings := new(protocol.SettingsProviderMock)
	utxosManagerMock := new(protocol.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)

	// Act
	actualTimestamp := blockchain.FirstBlockTimestamp()

	// Assert
	var expectedTimestamp int64 = 0
	test.Assert(t, actualTimestamp == expectedTimestamp, fmt.Sprintf("timestamp is %d whereas it should be %d", actualTimestamp, expectedTimestamp))
}

func Test_FirstBlockTimestamp_BlockchainIsNotEmpty_ReturnsFirstBlockTimestamp(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	settings := new(protocol.SettingsProviderMock)
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)

	// Act
	actualTimestamp := blockchain.FirstBlockTimestamp()

	// Assert
	expectedTimestamp := genesisTimestamp
	test.Assert(t, actualTimestamp == expectedTimestamp, fmt.Sprintf("timestamp is %d whereas it should be %d", actualTimestamp, expectedTimestamp))
}

func Test_LastBlockTimestamp_BlockchainIsEmpty_Returns0(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	settings := new(protocol.SettingsProviderMock)
	utxosManagerMock := new(protocol.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)

	// Act
	actualTimestamp := blockchain.LastBlockTimestamp()

	// Assert
	var expectedTimestamp int64 = 0
	test.Assert(t, actualTimestamp == expectedTimestamp, fmt.Sprintf("timestamp is %d whereas it should be %d", actualTimestamp, expectedTimestamp))
}

func Test_LastBlockTimestamp_BlockchainIsNotEmpty_ReturnsLastBlockTimestamp(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborsManagerMock := new(network.NeighborsManagerMock)
	settings := new(protocol.SettingsProviderMock)
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	var genesisTimestamp int64 = 0
	var expectedTimestamp int64 = 1
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(expectedTimestamp, nil, nil)

	// Act
	actualTimestamp := blockchain.LastBlockTimestamp()

	// Assert
	test.Assert(t, actualTimestamp == expectedTimestamp, fmt.Sprintf("timestamp is %d whereas it should be %d", actualTimestamp, expectedTimestamp))
}

func Test_Update_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	var validationTimestamp int64 = 11
	settings := new(protocol.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 2 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	now := 5 * validationTimestamp
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.ClearFunc = func() {}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(now-5*validationTimestamp, nil, nil)
	_ = blockchain.AddBlock(now-4*validationTimestamp, nil, nil)
	blocksBytes := blockchain.Blocks(0)
	var blocks []*ledger.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	genesisBlockHash := blocks[1].PreviousHash()
	block1 := ledger.NewRewardedBlock(genesisBlockHash, now-4*validationTimestamp)
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-3*validationTimestamp)
	hash2, _ := block2.Hash()
	block3 := ledger.NewRewardedBlock(hash2, now-2*validationTimestamp)
	hash3, _ := block3.Hash()
	block4 := ledger.NewRewardedBlock(hash3, now-validationTimestamp)
	neighborBlocks := []*ledger.Block{blocks[0], block1, block2, block3, block4}
	neighborBlocksBytes, _ := json.Marshal(neighborBlocks)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		return neighborBlocksBytes, nil
	}

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		blockchainReplacedMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTimestampIsInvalid_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

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
			neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
				block1 := ledger.NewRewardedBlock([32]byte{}, tt.args.firstBlockTimestamp)
				hash, _ := block1.Hash()
				block2 := ledger.NewRewardedBlock(hash, tt.args.secondBlockTimestamp)
				blocks := []*ledger.Block{block1, block2}
				blockBytes, _ := json.Marshal(blocks)
				return blockBytes, nil
			}

			// Act
			blockchain.Update(1)

			// Assert
			expectedMessages := []string{
				"neighbor block timestamp is invalid",
				blockchainKeptMessage,
			}
			test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
		})
	}
}

func Test_Update_NeighborNewBlockTimestampIsInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	var validationTimestamp int64 = 1
	now := validationTimestamp
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		block1 := ledger.NewRewardedBlock([32]byte{}, now)
		hash, _ := block1.Hash()
		block2 := ledger.NewRewardedBlock(hash, now+validationTimestamp)
		blocks := []*ledger.Block{block1, block2}
		blockBytes, _ := json.Marshal(blocks)
		return blockBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"neighbor block timestamp is in the future",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTransactionFeeIsNegative_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	address := test.Address
	var invalidTransactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var incomeLimit uint64 = 1
	genesisAmount := 2 * incomeLimit
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	genesisTransaction := block1.Transactions()[0]
	var genesisOutputIndex uint16 = 0
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, invalidTransactionFee, genesisOutputIndex, "A", privateKey, publicKey, now, genesisTransaction.Id(), genesisAmount, false)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 1)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, []string{address}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == invalidTransaction.Id() {
			return ledger.NewUtxo(invalidTransaction.Inputs()[0].InputInfo, invalidTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"transaction fee is negative",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTransactionFeeIsTooLow_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	address := test.Address
	var invalidTransactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	genesisTransaction := block1.Transactions()[0]
	var genesisOutputIndex uint16 = 0
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, invalidTransactionFee, genesisOutputIndex, "A", privateKey, publicKey, now, genesisTransaction.Id(), genesisAmount, false)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 1)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, []string{address}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == invalidTransaction.Id() {
			return ledger.NewUtxo(invalidTransaction.Inputs()[0].InputInfo, invalidTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"transaction fee is too low",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTransactionTimestampIsTooFarInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, transactionFee, genesisOutputIndex, "A", privateKey, publicKey, now+validationTimestamp, genesisTransaction.Id(), genesisAmount, false)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 0)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, []string{address}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == invalidTransaction.Id() {
			return ledger.NewUtxo(invalidTransaction.Inputs()[0].InputInfo, invalidTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		fmt.Sprintf("a neighbor block transaction timestamp is too far in the future: transaction timestamp: %d, id: %s", invalidTransaction.Timestamp(), invalidTransaction.Id()),
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTransactionTimestampIsTooOld_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, transactionFee, genesisOutputIndex, "A", privateKey, publicKey, now-validationTimestamp-1, genesisTransaction.Id(), genesisAmount, false)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 0)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, []string{address}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == invalidTransaction.Id() {
			return ledger.NewUtxo(invalidTransaction.Inputs()[0].InputInfo, invalidTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		fmt.Sprintf("a neighbor block transaction timestamp is too old: transaction timestamp: %d, id: %s", invalidTransaction.Timestamp(), invalidTransaction.Id()),
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTransactionInputSignatureIsInvalid_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, transactionFee, genesisOutputIndex, "A", privateKey2, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, false)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 0)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, []string{address}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == invalidTransaction.Id() {
			return ledger.NewUtxo(invalidTransaction.Inputs()[0].InputInfo, invalidTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"neighbor transaction is invalid: failed to verify signature of an input: signature is invalid",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTransactionInputPublicKeyIsInvalid_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	address := test.Address
	var transactionFee uint64 = 0
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	publicKey2 := encryption.NewPublicKey(privateKey2)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, transactionFee, genesisOutputIndex, "A", privateKey2, publicKey2, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, false)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 0)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, []string{address}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == invalidTransaction.Id() {
			return ledger.NewUtxo(invalidTransaction.Inputs()[0].InputInfo, invalidTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"neighbor transaction is invalid: output address does not derive from input public key",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborAddressIsNotRegistered_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	notRegisteredAddress := test.Address
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return errors.New("") }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := ledger.NewGenesisBlock(notRegisteredAddress, genesisAmount)
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(notRegisteredAddress, false, now, 0)
	transactions := []*ledger.Transaction{rewardTransaction}
	block3 := ledger.NewBlock(hash2, []string{notRegisteredAddress}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"failed to verify registered addresses",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborBlockYieldingOutputAddressIsRegistered_IsReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	address := test.Address
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, transactionFee, genesisOutputIndex, address, privateKey, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, true)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 0)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, nil, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == genesisTransaction.Id() {
			return ledger.NewUtxo(nil, genesisTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		blockchainReplacedMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborBlockYieldingOutputAddressHasBeenRecentlyAdded_IsReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return false }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	address := test.Address
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	addedAddress := test.Address2
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, transactionFee, genesisOutputIndex, addedAddress, privateKey, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, true)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 0)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, []string{addedAddress}, nil, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == genesisTransaction.Id() {
			return ledger.NewUtxo(nil, genesisTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		blockchainReplacedMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborBlockYieldingOutputIsNotRegistered_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return false }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	address := test.Address
	block1 := ledger.NewGenesisBlock(address, genesisAmount)
	removedAddress := test.Address2
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransactionRequestBytes := ledger.NewSignedTransactionRequest(genesisAmount, transactionFee, genesisOutputIndex, removedAddress, privateKey, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, true)
	var invalidTransactionRequest *ledger.TransactionRequest
	_ = json.Unmarshal(invalidTransactionRequestBytes, &invalidTransactionRequest)
	invalidTransaction := invalidTransactionRequest.Transaction()
	hash1, _ := block1.Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := ledger.NewRewardTransaction(address, false, now, 0)
	transactions := []*ledger.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := ledger.NewBlock(hash2, nil, []string{removedAddress}, now, transactions)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*ledger.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		if input.TransactionId() == genesisTransaction.Id() {
			return ledger.NewUtxo(nil, genesisTransaction.Outputs()[0], now), nil
		} else {
			return ledger.NewUtxo(nil, rewardTransaction.Outputs()[0], now), nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"neighbor block transaction yielding output address is not registered",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborValidatorIsNotTheOldest_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 1 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	rewardTransaction1, _ := ledger.NewRewardTransaction(test.Address, false, now-2*validationTimestamp, 0)
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		return ledger.NewUtxo(nil, rewardTransaction1.Outputs()[0], now), nil
	}
	transactionsBytes1, _ := json.Marshal([]*ledger.Transaction{rewardTransaction1})
	_ = blockchain.AddBlock(now-2*validationTimestamp, transactionsBytes1, nil)
	blocksBytes := blockchain.Blocks(0)
	var blocks []*ledger.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	rewardTransaction2, _ := ledger.NewRewardTransaction(test.Address, false, now-validationTimestamp, 0)
	transactionsBytes2, _ := json.Marshal([]*ledger.Transaction{rewardTransaction2})
	_ = blockchain.AddBlock(now-validationTimestamp, transactionsBytes2, nil)
	rewardTransaction3, _ := ledger.NewRewardTransaction(test.Address, false, now, 0)
	transactionsBytes3, _ := json.Marshal([]*ledger.Transaction{rewardTransaction3})
	_ = blockchain.AddBlock(now, transactionsBytes3, nil)
	hash1, _ := blocks[0].Hash()
	block2 := ledger.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	block3 := ledger.NewRewardedBlock(hash2, now)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		neighborBlocks := []*ledger.Block{block3}
		neighborBlocksBytes, _ := json.Marshal(neighborBlocks)
		return neighborBlocksBytes, nil
	}

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborValidatorIsTheOldest_IsReplaced(t *testing.T) {
	// Arrange
	registryMock := new(protocol.AddressesManagerMock)
	registryMock.CopyFunc = func() protocol.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	neighborMock := new(presentation.NeighborCallerMock)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	neighborsManagerMock := new(network.NeighborsManagerMock)
	neighborsManagerMock.NeighborsFunc = func() []presentation.NeighborCaller {
		return []presentation.NeighborCaller{neighborMock}
	}
	settings := new(protocol.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 2 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(protocol.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() protocol.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]byte, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, neighborsManagerMock, utxosManagerMock, logger)
	rewardTransaction1, _ := ledger.NewRewardTransaction(test.Address, false, now-2*validationTimestamp, 0)
	transactionsBytes1, _ := json.Marshal([]*ledger.Transaction{rewardTransaction1})
	utxosManagerMock.UtxoFunc = func(input protocol.InputInfoProvider) (protocol.UtxoInfoProvider, error) {
		return ledger.NewUtxo(nil, rewardTransaction1.Outputs()[0], now), nil
	}
	_ = blockchain.AddBlock(now-2*validationTimestamp, transactionsBytes1, nil)
	rewardTransaction2, _ := ledger.NewRewardTransaction(test.Address, false, now-validationTimestamp, 0)
	transactionsBytes2, _ := json.Marshal([]*ledger.Transaction{rewardTransaction2})
	_ = blockchain.AddBlock(now-validationTimestamp, transactionsBytes2, nil)
	blocksBytes := blockchain.Blocks(0)
	var blocks []*ledger.Block
	_ = json.Unmarshal(blocksBytes, &blocks)
	rewardTransaction3, _ := ledger.NewRewardTransaction(test.Address, false, now, 0)
	transactionsBytes3, _ := json.Marshal([]*ledger.Transaction{rewardTransaction3})
	_ = blockchain.AddBlock(now, transactionsBytes3, nil)
	hash2, _ := blocks[1].Hash()
	block3 := ledger.NewRewardedBlock(hash2, now)
	neighborMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		neighborBlocks := []*ledger.Block{block3}
		neighborBlocksBytes, _ := json.Marshal(neighborBlocks)
		return neighborBlocksBytes, nil
	}

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		blockchainReplacedMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}
