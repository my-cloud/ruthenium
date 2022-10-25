package protocol

import (
	"github.com/my-cloud/ruthenium/src/api/node"
	"github.com/my-cloud/ruthenium/src/api/node/network"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_Verify_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	registrable := new(RegistrableMock)
	registrable.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	watch := clock.NewWatch()
	logger := log.NewLogger(log.Fatal)
	requestable := new(RequestableMock)
	requestable.GetBlocksFunc = func() ([]*node.BlockResponse, error) {
		blockResponse1 := &node.BlockResponse{
			Timestamp:           0,
			PreviousHash:        [32]byte{},
			Transactions:        []*node.TransactionResponse{protocol.NewRewardTransaction("recipient", 0, 0).GetResponse()},
			RegisteredAddresses: nil,
		}
		block1, _ := protocol.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := &node.BlockResponse{
			Timestamp:           0,
			PreviousHash:        hash,
			Transactions:        []*node.TransactionResponse{protocol.NewRewardTransaction("recipient", 0, 0).GetResponse()},
			RegisteredAddresses: nil,
		}
		return []*node.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	requestable.TargetFunc = func() string {
		return "requestable"
	}
	synchronizable := new(SynchronizableMock)
	synchronizable.NeighborsFunc = func() []network.Requestable {
		return []network.Requestable{requestable}
	}
	blockchain := protocol.NewBlockchain(registrable, 0, watch, synchronizable, logger)
	blockchain.Verify()
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, isReplaced, "blockchain is not replaced whereas it should be")
}
