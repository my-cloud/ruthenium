package protocol

import (
	"github.com/my-cloud/ruthenium/src/api/node/network"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

//func Test_IsEmpty_Empty_ReturnsTrue(t *testing.T) {
//	blockchain := protocol.NewBlockchain(0, nil, nil)
//	isEmpty := blockchain.IsEmpty()
//	test.Assert(t, isEmpty, "blockchain should be empty")
//}
//
//func Test_IsEmpty_BlockAdded_ReturnsFalse(t *testing.T) {
//	blockchain := protocol.NewBlockchain(0, nil, nil)
//	blockchain.AddBlock(protocol.NewBlock(0, [32]byte{}, nil, nil))
//	isEmpty := blockchain.IsEmpty()
//	test.Assert(t, !isEmpty, "blockchain should not be empty")
//}
//
//func Test_BlockResponses_Empty_ReturnsNil(t *testing.T) {
//	blockchain := protocol.NewBlockchain(0, nil, nil)
//	blockResponses := blockchain.Blocks()
//	test.Assert(t, blockResponses == nil, "block responses should be nil")
//}
//
//func Test_CalculateTotalAmount__ReturnsTotalAmount(t *testing.T) {
//	blockchain := protocol.NewBlockchain(0, nil, nil)
//	transaction := protocol.NewTransaction("recipient", "sender", nil, 0, 10)
//	block := protocol.NewBlock(0, [32]byte{}, []*protocol.Transaction{transaction}, nil)
//	blockchain.AddBlock(block)
//	actualTotalAmount := blockchain.CalculateTotalAmount(0, "recipient")
//	var expectedTotalAmount uint64 = 10
//	test.Assert(t, actualTotalAmount == expectedTotalAmount, fmt.Sprintf("Wrong total amount. Expected: %d - Actual: %d", expectedTotalAmount, actualTotalAmount))
//}
func Test_Verify_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	registrable := NewRegistrableMock()
	IsRegisteredMock = func(address string) (bool, error) { return true, nil }
	watch := clock.NewWatch()
	logger := log.NewLogger(log.Fatal)
	blockchain := protocol.NewBlockchain(registrable, 0, watch, logger)
	requestable := NewRequestableMock()
	GetBlocksMock = func() ([]*network.BlockResponse, error) {
		blockResponse1 := &network.BlockResponse{
			Timestamp:           0,
			PreviousHash:        [32]byte{},
			Transactions:        []*network.TransactionResponse{protocol.NewRewardTransaction("recipient", 0, 0).GetResponse()},
			RegisteredAddresses: nil,
		}
		block1, _ := protocol.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := &network.BlockResponse{
			Timestamp:           0,
			PreviousHash:        hash,
			Transactions:        []*network.TransactionResponse{protocol.NewRewardTransaction("recipient", 0, 0).GetResponse()},
			RegisteredAddresses: nil,
		}
		return []*network.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	TargetMock = func() string {
		return "requestable"
	}
	blockchain.Verify([]network.Requestable{requestable})
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, isReplaced, "blockchain is not replaced whereas it should be")
}
