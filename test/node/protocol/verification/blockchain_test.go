package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"github.com/my-cloud/ruthenium/test/node/clock/clocktest"
	"github.com/my-cloud/ruthenium/test/node/network/networktest"
	"github.com/my-cloud/ruthenium/test/node/protocol/protocoltest"
	"strings"
	"testing"
	"time"
)

const (
	blockchainReplacedMessage = "verification done: blockchain replaced"
	blockchainKeptMessage     = "verification done: blockchain kept"
)

func Test_AddBlock_ValidParameters_NoErrorLogged(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

	// Act
	blockchain.AddBlock(0, nil, nil)

	// Assert
	test.Assert(t, len(logger.ErrorCalls()) == 0, "logger has been called whereas it should not")
}

func Test_Blocks_ValidParameters_NoErrorLogged(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	logger := logtest.NewLoggerMock()
	synchronizer := new(networktest.SynchronizerMock)
	genesisTransaction := validation.NewRewardTransaction("", 0, 0)
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

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
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

	// Act
	amount := blockchain.CalculateTotalAmount(1, genesisTransaction.RecipientAddress)

	// Assert
	test.Assert(t, amount == genesisTransaction.Value, "calculated amount is not the genesis amount whereas it should be")
}

func Test_Update_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	// Arrange
	registry := new(protocoltest.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewRewardedBlockResponse([32]byte{}, 0)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := protocoltest.NewRewardedBlockResponse(hash, 1)
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
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(watchMock.Now().UnixNano())

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
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 2) }
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
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

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
				block1, _ := verification.NewBlockFromResponse(blockResponse1)
				hash, _ := block1.Hash()
				blockResponse2 := protocoltest.NewRewardedBlockResponse(hash, tt.args.secondBlockTimestamp)
				return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
			}

			// Act
			blockchain.Update(watchMock.Now().UnixNano())

			// Assert
			var isKept bool
			var isExplicitMessageLogged bool
			for _, call := range logger.DebugCalls() {
				expectedMessage := fmt.Sprintf("neighbor block timestamp is invalid")
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
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewRewardedBlockResponse([32]byte{}, 1)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
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
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(watchMock.Now().UnixNano())

	// Assert
	var isKept bool
	var isExplicitMessageLogged bool
	for _, call := range logger.DebugCalls() {
		expectedMessage := fmt.Sprintf("neighbor block timestamp is in the future")
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
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	wallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	address := wallet.Address()
	serverTransaction := server.NewTransaction("A", address, wallet.PublicKey(), 3, 1)
	_ = serverTransaction.Sign(wallet.PrivateKey())
	transactionRequest := serverTransaction.GetRequest()
	transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
	transactionResponse := transaction.GetResponse()
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		var block2Timestamp int64 = 1
		transactions := []*network.TransactionResponse{
			transactionResponse,
			validation.NewRewardTransaction(address, block2Timestamp, 0),
		}
		var registeredAddresses []string
		registeredAddresses = append(registeredAddresses, address)
		blockResponse2 := verification.NewBlockResponse(block2Timestamp, hash, transactions, registeredAddresses)
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
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(watchMock.Now().UnixNano())

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
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	watchMock := new(clocktest.WatchMock)
	watchMock.NowFunc = func() time.Time { return time.Unix(0, 2) }
	logger := logtest.NewLoggerMock()
	neighborMock := new(networktest.NeighborMock)
	wallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
	address := wallet.Address()
	serverTransaction := server.NewTransaction("A", address, wallet.PublicKey(), 0, 1)
	_ = serverTransaction.Sign(wallet.PrivateKey())
	transactionRequest := serverTransaction.GetRequest()
	transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
	transactionResponse := transaction.GetResponse()
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := protocoltest.NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash1, _ := block1.Hash()
		var block2Timestamp int64 = 1
		transactions := []*network.TransactionResponse{
			transactionResponse,
			validation.NewRewardTransaction(address, block2Timestamp, 0),
		}
		var registeredAddresses []string
		registeredAddresses = append(registeredAddresses, address)
		blockResponse2 := verification.NewBlockResponse(block2Timestamp, hash1, transactions, registeredAddresses)
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
	blockchain := verification.NewBlockchain(genesisTransaction, registry, 1, synchronizer, logger)

	// Act
	blockchain.Update(watchMock.Now().UnixNano())

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
