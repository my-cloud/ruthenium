package verification

import (
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	network2 "github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/mock"
	"testing"
	"time"
)

func Test_Verify_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	// Arrange
	registry := new(mock.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(mock.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network2.BlockResponse, error) {
		blockResponse1 := mock.NewRewardedBlockResponse([32]byte{}, 0)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := mock.NewRewardedBlockResponse(hash, 1)
		return []*network2.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(mock.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network2.Neighbor {
		return []network2.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, synchronizer, logger)

	// Act
	blockchain.Verify(timeMock.Now().UnixNano())

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, isReplaced, "blockchain is not replaced whereas it should be")
}

func Test_Verify_NeighborNewBlockTimestampIsInvalid_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(mock.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 2) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(mock.NeighborMock)
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(mock.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network2.Neighbor {
		return []network2.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, synchronizer, logger)

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
			neighborMock.GetBlocksFunc = func() ([]*network2.BlockResponse, error) {
				blockResponse1 := mock.NewRewardedBlockResponse([32]byte{}, tt.args.firstBlockTimestamp)
				block1, _ := verification.NewBlockFromResponse(blockResponse1)
				hash, _ := block1.Hash()
				blockResponse2 := mock.NewRewardedBlockResponse(hash, tt.args.secondBlockTimestamp)
				return []*network2.BlockResponse{blockResponse1, blockResponse2}, nil
			}

			// Act
			blockchain.Verify(timeMock.Now().UnixNano())

			// Assert
			if blockchain.IsReplaced() {
				t.Errorf("blockchain is replaced whereas it should not")
			}
		})
	}
}

func Test_Verify_NeighborNewBlockTimestampIsInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(mock.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(mock.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network2.BlockResponse, error) {
		blockResponse1 := mock.NewRewardedBlockResponse([32]byte{}, 1)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := mock.NewRewardedBlockResponse(hash, 2)
		return []*network2.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(mock.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network2.Neighbor {
		return []network2.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, synchronizer, logger)

	// Act
	blockchain.Verify(timeMock.Now().UnixNano())

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, !isReplaced, "blockchain is replaced whereas it should not")
}

func Test_Verify_NeighborNewBlockTransactionTimestampIsTooFarInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(mock.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(mock.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network2.BlockResponse, error) {
		wallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
		address := wallet.Address()
		blockResponse1 := mock.NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		var block2Timestamp int64 = 1
		serverTransaction := server.NewTransaction("A", wallet.Address(), wallet.PublicKey(), 3, 1)
		_ = serverTransaction.Sign(wallet.PrivateKey())
		transactionRequest := serverTransaction.GetRequest()
		transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
		transactions := []*network2.TransactionResponse{validation.NewRewardTransaction(address, block2Timestamp, 0), transaction.GetResponse()}
		var registeredAddresses []string
		registeredAddresses = append(registeredAddresses, address)
		blockResponse2 := mock.NewBlockResponse(block2Timestamp, hash, transactions, registeredAddresses)
		return []*network2.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(mock.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network2.Neighbor {
		return []network2.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, synchronizer, logger)

	// Act
	blockchain.Verify(timeMock.Now().UnixNano())

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, !isReplaced, "blockchain is replaced whereas it should not")
}

func Test_Verify_NeighborNewBlockTransactionTimestampIsTooOld_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(mock.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(mock.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 2) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(mock.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network2.BlockResponse, error) {
		wallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
		address := wallet.Address()
		blockResponse1 := mock.NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash1, _ := block1.Hash()
		blockResponse2 := mock.NewRewardedBlockResponse(hash1, 1)
		block2, _ := verification.NewBlockFromResponse(blockResponse2)
		hash2, _ := block2.Hash()
		var block3Timestamp int64 = 2
		serverTransaction := server.NewTransaction("A", wallet.Address(), wallet.PublicKey(), 0, 1)
		_ = serverTransaction.Sign(wallet.PrivateKey())
		transactionRequest := serverTransaction.GetRequest()
		transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
		transactions := []*network2.TransactionResponse{validation.NewRewardTransaction(address, block3Timestamp, 0), transaction.GetResponse()}
		var registeredAddresses []string
		registeredAddresses = append(registeredAddresses, address)
		blockResponse3 := mock.NewBlockResponse(block3Timestamp, hash2, transactions, registeredAddresses)
		return []*network2.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(mock.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network2.Neighbor {
		return []network2.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, synchronizer, logger)

	// Act
	blockchain.Verify(timeMock.Now().UnixNano())

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, !isReplaced, "blockchain is replaced whereas it should not")
}
