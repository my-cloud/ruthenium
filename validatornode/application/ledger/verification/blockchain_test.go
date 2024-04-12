package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

const (
	blockchainReplacedMessage = "verification done: blockchain replaced"
	blockchainKeptMessage     = "verification done: blockchain kept"
)

func Test_AddBlock_ValidParameters_NoErrorReturned(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	settings := new(ledger.SettingsProviderMock)
	utxosManagerMock := new(ledger.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)

	// Act
	err := blockchain.AddBlock(0, nil, nil)

	// Assert
	test.Assert(t, err == nil, "error is returned whereas it should not")
}

func Test_Blocks_BlocksCountLimitSetToZero_ReturnsEmptyArray(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	settings := new(ledger.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 0 }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)

	// Act
	blocks := blockchain.Blocks(0)

	// Assert
	test.Assert(t, len(blocks) == 0, "blocks should be empty")
}

func Test_Blocks_BlocksCountLimitSetToOne_ReturnsOneBlock(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	var expectedBlocksCount uint64 = 1
	settings := new(ledger.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return expectedBlocksCount }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
	var validationInterval int64 = 1
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, nil, nil)

	// Act
	blocks := blockchain.Blocks(0)

	// Assert
	actualBlocksCount := uint64(len(blocks))
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_Blocks_BlocksCountLimitSetToTwo_ReturnsTwoBlocks(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	var expectedBlocksCount uint64 = 2
	settings := new(ledger.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return expectedBlocksCount }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
	var validationInterval int64 = 1
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, nil, nil)

	// Act
	blocks := blockchain.Blocks(0)

	// Assert
	actualBlocksCount := uint64(len(blocks))
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_Blocks_StartingBlockHeightGreaterThanBlocksLength_ReturnsEmptyArray(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	var blocksCount uint64 = 1
	settings := new(ledger.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return blocksCount }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
	var genesisTimestamp int64 = 0
	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)

	// Act
	blocks := blockchain.Blocks(1)

	// Assert
	expectedBlocksCount := 0
	actualBlocksCount := len(blocks)
	test.Assert(t, actualBlocksCount == expectedBlocksCount, fmt.Sprintf("blocks count is %d whereas it should be %d", actualBlocksCount, expectedBlocksCount))
}

func Test_FirstBlockTimestamp_BlockchainIsEmpty_Returns0(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	settings := new(ledger.SettingsProviderMock)
	utxosManagerMock := new(ledger.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)

	// Act
	actualTimestamp := blockchain.FirstBlockTimestamp()

	// Assert
	var expectedTimestamp int64 = 0
	test.Assert(t, actualTimestamp == expectedTimestamp, fmt.Sprintf("timestamp is %d whereas it should be %d", actualTimestamp, expectedTimestamp))
}

func Test_FirstBlockTimestamp_BlockchainIsNotEmpty_ReturnsFirstBlockTimestamp(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	settings := new(ledger.SettingsProviderMock)
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	settings := new(ledger.SettingsProviderMock)
	utxosManagerMock := new(ledger.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)

	// Act
	actualTimestamp := blockchain.LastBlockTimestamp()

	// Assert
	var expectedTimestamp int64 = 0
	test.Assert(t, actualTimestamp == expectedTimestamp, fmt.Sprintf("timestamp is %d whereas it should be %d", actualTimestamp, expectedTimestamp))
}

func Test_LastBlockTimestamp_BlockchainIsNotEmpty_ReturnsLastBlockTimestamp(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	sendersManagerMock := new(network.SendersManagerMock)
	settings := new(ledger.SettingsProviderMock)
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	var validationTimestamp int64 = 11
	settings := new(ledger.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 2 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	now := 5 * validationTimestamp
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.ClearFunc = func() {}
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(now-5*validationTimestamp, nil, nil)
	_ = blockchain.AddBlock(now-4*validationTimestamp, nil, nil)
	blocks := blockchain.Blocks(0)
	genesisBlockHash := blocks[1].PreviousHash()
	block1 := protocol.NewRewardedBlock(genesisBlockHash, now-4*validationTimestamp)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-3*validationTimestamp)
	hash2, _ := block2.Hash()
	block3 := protocol.NewRewardedBlock(hash2, now-2*validationTimestamp)
	hash3, _ := block3.Hash()
	block4 := protocol.NewRewardedBlock(hash3, now-validationTimestamp)
	neighborBlocks := []*protocol.Block{blocks[0], block1, block2, block3, block4}
	neighborBlocksBytes, _ := json.Marshal(neighborBlocks)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return 1 }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
			senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
				block1 := protocol.NewRewardedBlock([32]byte{}, tt.args.firstBlockTimestamp)
				hash, _ := block1.Hash()
				block2 := protocol.NewRewardedBlock(hash, tt.args.secondBlockTimestamp)
				blocks := []*protocol.Block{block1, block2}
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	var validationTimestamp int64 = 1
	now := validationTimestamp
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		block1 := protocol.NewRewardedBlock([32]byte{}, now)
		hash, _ := block1.Hash()
		block2 := protocol.NewRewardedBlock(hash, now+validationTimestamp)
		blocks := []*protocol.Block{block1, block2}
		blockBytes, _ := json.Marshal(blocks)
		return blockBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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

func Test_Update_NeighborNewBlockTransactionFeeCalculationFails_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	address := test.Address
	var invalidTransactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var incomeLimit uint64 = 1
	genesisAmount := 2 * incomeLimit
	block1 := protocol.NewGenesisBlock(address, genesisAmount)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	genesisTransaction := block1.Transactions()[0]
	var genesisOutputIndex uint16 = 0
	invalidTransaction := protocol.NewSignedTransaction(genesisAmount, invalidTransactionFee, genesisOutputIndex, "A", privateKey, publicKey, now, genesisTransaction.Id(), genesisAmount, false)
	rewardTransaction, _ := protocol.NewRewardTransaction(address, false, now, 1)
	transactions := []*protocol.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := protocol.NewBlock(hash2, []string{address}, nil, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 0 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) {
		if transaction.Id() == invalidTransaction.Id() {
			return 0, errors.New("")
		} else {
			return 0, nil
		}
	}
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
	_ = blockchain.AddBlock(0, nil, nil)

	// Act
	blockchain.Update(now)

	// Assert
	expectedMessages := []string{
		"failed to verify a neighbor block transaction fee",
		blockchainKeptMessage,
	}
	test.AssertThatMessageIsLogged(t, expectedMessages, logger.DebugCalls())
}

func Test_Update_NeighborNewBlockTransactionTimestampIsTooFarInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := protocol.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransaction := protocol.NewSignedTransaction(genesisAmount, transactionFee, genesisOutputIndex, "A", privateKey, publicKey, now+validationTimestamp, genesisTransaction.Id(), genesisAmount, false)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := protocol.NewRewardTransaction(address, false, now, 0)
	transactions := []*protocol.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := protocol.NewBlock(hash2, []string{address}, nil, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := protocol.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransaction := protocol.NewSignedTransaction(genesisAmount, transactionFee, genesisOutputIndex, "A", privateKey, publicKey, now-validationTimestamp-1, genesisTransaction.Id(), genesisAmount, false)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := protocol.NewRewardTransaction(address, false, now, 0)
	transactions := []*protocol.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := protocol.NewBlock(hash2, []string{address}, nil, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	address := test.Address
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	privateKey2, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey2)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := protocol.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransaction := protocol.NewSignedTransaction(genesisAmount, transactionFee, genesisOutputIndex, "A", privateKey2, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, false)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := protocol.NewRewardTransaction(address, false, now, 0)
	transactions := []*protocol.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := protocol.NewBlock(hash2, []string{address}, nil, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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

func Test_Update_NeighborAddressIsNotRegistered_IsNotReplaced(t *testing.T) {
	// Arrange
	registryMock := new(ledger.AddressesManagerMock)
	notRegisteredAddress := test.Address
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return errors.New("") }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	block1 := protocol.NewGenesisBlock(notRegisteredAddress, genesisAmount)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := protocol.NewRewardTransaction(notRegisteredAddress, false, now, 0)
	transactions := []*protocol.Transaction{rewardTransaction}
	block3 := protocol.NewBlock(hash2, []string{notRegisteredAddress}, nil, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	address := test.Address
	block1 := protocol.NewGenesisBlock(address, genesisAmount)
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransaction := protocol.NewSignedTransaction(genesisAmount, transactionFee, genesisOutputIndex, address, privateKey, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, true)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := protocol.NewRewardTransaction(address, false, now, 0)
	transactions := []*protocol.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := protocol.NewBlock(hash2, nil, nil, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return false }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	address := test.Address
	block1 := protocol.NewGenesisBlock(address, genesisAmount)
	addedAddress := test.Address2
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransaction := protocol.NewSignedTransaction(genesisAmount, transactionFee, genesisOutputIndex, addedAddress, privateKey, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, true)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := protocol.NewRewardTransaction(address, false, now, 0)
	transactions := []*protocol.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := protocol.NewBlock(hash2, []string{addedAddress}, nil, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.ClearFunc = func() {}
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return false }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	var transactionFee uint64 = 0
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	var genesisAmount uint64 = 1
	address := test.Address
	block1 := protocol.NewGenesisBlock(address, genesisAmount)
	removedAddress := test.Address2
	var genesisOutputIndex uint16 = 0
	genesisTransaction := block1.Transactions()[0]
	invalidTransaction := protocol.NewSignedTransaction(genesisAmount, transactionFee, genesisOutputIndex, removedAddress, privateKey, publicKey, now-validationTimestamp, genesisTransaction.Id(), genesisAmount, true)
	hash1, _ := block1.Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	rewardTransaction, _ := protocol.NewRewardTransaction(address, false, now, 0)
	transactions := []*protocol.Transaction{
		invalidTransaction,
		rewardTransaction,
	}
	block3 := protocol.NewBlock(hash2, nil, []string{removedAddress}, now, transactions)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		blocks := []*protocol.Block{block1, block2, block3}
		blocksBytes, _ := json.Marshal(blocks)
		return blocksBytes, nil
	}
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return transactionFee }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.ClearFunc = func() {}
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 1 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
	rewardTransaction1, _ := protocol.NewRewardTransaction(test.Address, false, now-2*validationTimestamp, 0)
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	_ = blockchain.AddBlock(now-2*validationTimestamp, []*protocol.Transaction{rewardTransaction1}, nil)
	blocks := blockchain.Blocks(0)
	rewardTransaction2, _ := protocol.NewRewardTransaction(test.Address, false, now-validationTimestamp, 0)
	_ = blockchain.AddBlock(now-validationTimestamp, []*protocol.Transaction{rewardTransaction2}, nil)
	rewardTransaction3, _ := protocol.NewRewardTransaction(test.Address, false, now, 0)
	_ = blockchain.AddBlock(now, []*protocol.Transaction{rewardTransaction3}, nil)
	hash1, _ := blocks[0].Hash()
	block2 := protocol.NewRewardedBlock(hash1, now-validationTimestamp)
	hash2, _ := block2.Hash()
	block3 := protocol.NewRewardedBlock(hash2, now)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		neighborBlocks := []*protocol.Block{block3}
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
	registryMock := new(ledger.AddressesManagerMock)
	registryMock.CopyFunc = func() ledger.AddressesManager { return registryMock }
	registryMock.FilterFunc = func([]string) []string { return nil }
	registryMock.IsRegisteredFunc = func(string) bool { return true }
	registryMock.RemovedAddressesFunc = func() []string { return nil }
	registryMock.VerifyFunc = func([]string, []string) error { return nil }
	registryMock.UpdateFunc = func([]string, []string) {}
	logger := log.NewLoggerMock()
	senderMock := new(network.SenderMock)
	var validationTimestamp int64 = 1
	now := 2 * validationTimestamp
	senderMock.TargetFunc = func() string {
		return "neighbor"
	}
	sendersManagerMock := new(network.SendersManagerMock)
	sendersManagerMock.SendersFunc = func() []network.Sender {
		return []network.Sender{senderMock}
	}
	settings := new(ledger.SettingsProviderMock)
	settings.BlocksCountLimitFunc = func() uint64 { return 2 }
	settings.IncomeBaseFunc = func() uint64 { return 0 }
	settings.IncomeLimitFunc = func() uint64 { return 1 }
	settings.HalfLifeInNanosecondsFunc = func() float64 { return 0 }
	settings.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	settings.ValidationTimestampFunc = func() int64 { return validationTimestamp }
	settings.ValidationTimeoutFunc = func() time.Duration { return time.Second }
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.CopyFunc = func() ledger.UtxosManager { return utxosManagerMock }
	utxosManagerMock.UpdateUtxosFunc = func([]*protocol.Transaction, int64) error { return nil }
	blockchain := NewBlockchain(registryMock, settings, sendersManagerMock, utxosManagerMock, logger)
	rewardTransaction1, _ := protocol.NewRewardTransaction(test.Address, false, now-2*validationTimestamp, 0)
	utxosManagerMock.CalculateFeeFunc = func(transaction *protocol.Transaction, timestamp int64) (uint64, error) { return 0, nil }
	_ = blockchain.AddBlock(now-2*validationTimestamp, []*protocol.Transaction{rewardTransaction1}, nil)
	rewardTransaction2, _ := protocol.NewRewardTransaction(test.Address, false, now-validationTimestamp, 0)
	_ = blockchain.AddBlock(now-validationTimestamp, []*protocol.Transaction{rewardTransaction2}, nil)
	blocks := blockchain.Blocks(0)
	rewardTransaction3, _ := protocol.NewRewardTransaction(test.Address, false, now, 0)
	_ = blockchain.AddBlock(now, []*protocol.Transaction{rewardTransaction3}, nil)
	hash2, _ := blocks[1].Hash()
	block3 := protocol.NewRewardedBlock(hash2, now)
	senderMock.GetBlocksFunc = func(uint64) ([]byte, error) {
		neighborBlocks := []*protocol.Block{block3}
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
