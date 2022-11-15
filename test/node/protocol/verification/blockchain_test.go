package verification

import (
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/protocol"
	"testing"
	"time"
)

func Test_Verify_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	// Arrange
	registry := new(protocol.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(protocol.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(protocol.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := NewRewardedBlockResponse([32]byte{}, 0)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := NewRewardedBlockResponse(hash, 1)
		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(protocol.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

	// Act
	blockchain.Verify()

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, isReplaced, "blockchain is not replaced whereas it should be")
}

func Test_Verify_NeighborNewBlockTimestampIsInvalid_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocol.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(protocol.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 2) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(protocol.NeighborMock)
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(protocol.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

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
				blockResponse1 := NewRewardedBlockResponse([32]byte{}, tt.args.firstBlockTimestamp)
				block1, _ := verification.NewBlockFromResponse(blockResponse1)
				hash, _ := block1.Hash()
				blockResponse2 := NewRewardedBlockResponse(hash, tt.args.secondBlockTimestamp)
				return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
			}

			// Act
			blockchain.Verify()

			// Assert
			if blockchain.IsReplaced() {
				t.Errorf("blockchain is replaced whereas it should not")
			}
		})
	}
}

func Test_Verify_NeighborNewBlockTimestampIsInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocol.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(protocol.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(protocol.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		blockResponse1 := NewRewardedBlockResponse([32]byte{}, 1)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := NewRewardedBlockResponse(hash, 2)
		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(protocol.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

	// Act
	blockchain.Verify()

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, !isReplaced, "blockchain is replaced whereas it should not")
}

func Test_Verify_NeighborNewBlockTransactionTimestampIsTooFarInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocol.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(protocol.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(protocol.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		wallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
		address := wallet.Address()
		blockResponse1 := NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		var block2Timestamp int64 = 1
		serverTransaction := server.NewTransaction("A", wallet.Address(), wallet.PublicKey(), 3, 1)
		_ = serverTransaction.Sign(wallet.PrivateKey())
		transactionRequest := serverTransaction.GetRequest()
		transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
		transactions := []*validation.Transaction{validation.NewRewardTransaction(address, block2Timestamp, 0), transaction}
		blockResponse2 := NewBlockResponse(block2Timestamp, hash, transactions...)
		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(protocol.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

	// Act
	blockchain.Verify()

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, !isReplaced, "blockchain is replaced whereas it should not")
}

func Test_Verify_NeighborNewBlockTransactionTimestampIsTooOld_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(protocol.RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(protocol.TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 2) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(protocol.NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*network.BlockResponse, error) {
		wallet, _ := encryption.DecodeWallet(test.Mnemonic1, test.DerivationPath, "", "")
		address := wallet.Address()
		blockResponse1 := NewGenesisBlockResponse(address)
		block1, _ := verification.NewBlockFromResponse(blockResponse1)
		hash1, _ := block1.Hash()
		blockResponse2 := NewRewardedBlockResponse(hash1, 1)
		block2, _ := verification.NewBlockFromResponse(blockResponse2)
		hash2, _ := block2.Hash()
		var block3Timestamp int64 = 2
		serverTransaction := server.NewTransaction("A", wallet.Address(), wallet.PublicKey(), 0, 1)
		_ = serverTransaction.Sign(wallet.PrivateKey())
		transactionRequest := serverTransaction.GetRequest()
		transaction, _ := validation.NewTransactionFromRequest(&transactionRequest)
		transactions := []*validation.Transaction{validation.NewRewardTransaction(address, block3Timestamp, 0), transaction}
		blockResponse3 := NewBlockResponse(block3Timestamp, hash2, transactions...)
		return []*network.BlockResponse{blockResponse1, blockResponse2, blockResponse3}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(protocol.SynchronizerMock)
	synchronizer.NeighborsFunc = func() []network.Neighbor {
		return []network.Neighbor{neighborMock}
	}
	blockchain := verification.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

	// Act
	blockchain.Verify()

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, !isReplaced, "blockchain is replaced whereas it should not")
}
