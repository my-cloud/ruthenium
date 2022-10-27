package protocol

import (
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_Verify_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	// Arrange
	registry := new(RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	watch := clock.NewWatch()
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*neighborhood.BlockResponse, error) {
		blockResponse1 := &neighborhood.BlockResponse{
			Timestamp:           0,
			PreviousHash:        [32]byte{},
			Transactions:        []*neighborhood.TransactionResponse{protocol.NewRewardTransaction("recipient", 0, 0).GetResponse()},
			RegisteredAddresses: nil,
		}
		block1, _ := protocol.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := &neighborhood.BlockResponse{
			Timestamp:           0,
			PreviousHash:        hash,
			Transactions:        []*neighborhood.TransactionResponse{protocol.NewRewardTransaction("recipient", 0, 0).GetResponse()},
			RegisteredAddresses: nil,
		}
		return []*neighborhood.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(SynchronizerMock)
	synchronizer.NeighborsFunc = func() []neighborhood.Neighbor {
		return []neighborhood.Neighbor{neighborMock}
	}
	blockchain := protocol.NewBlockchain(registry, 0, watch, synchronizer, logger)

	// Act
	blockchain.Verify()

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, isReplaced, "blockchain is not replaced whereas it should be")
}
