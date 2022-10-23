package protocol

import (
	"github.com/my-cloud/ruthenium/src/api/node"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/humanity"
	"github.com/my-cloud/ruthenium/test/node/neighborhood"
	"testing"
)

func Test_Verify_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	registrable := humanity.NewRegistrableMock()
	humanity.IsRegisteredMock = func(address string) (bool, error) { return true, nil }
	watch := clock.NewWatch()
	logger := log.NewLogger(log.Fatal)
	blockchain := protocol.NewBlockchain(registrable, 0, watch, logger)
	requestable := neighborhood.NewRequestableMock()
	neighborhood.GetBlocksMock = func() ([]*node.BlockResponse, error) {
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
	neighborhood.TargetMock = func() string {
		return "requestable"
	}
	blockchain.Verify([]node.Requestable{requestable})
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, isReplaced, "blockchain is not replaced whereas it should be")
}
