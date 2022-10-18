package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_IsValidatorRegistered_Invalid_ReturnsFalse(t *testing.T) {
	senderAddress := protocol.RewardSenderAddress
	recipientAddress := "0x0000000000000000000000000000000000000001"
	transactions := []*protocol.Transaction{protocol.NewTransaction(recipientAddress, senderAddress, nil, 0, 1)}
	block := protocol.NewBlock(0, [32]byte{}, transactions, nil)
	isRegistered, _ := block.IsValidatorRegistered()
	test.Assert(t, !isRegistered, "proof of humanity is valid whereas it should not")
}

func Test_IsValidatorRegistered_Valid_ReturnsTrue(t *testing.T) {
	senderAddress := protocol.RewardSenderAddress
	recipientAddress := "0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a"
	transactions := []*protocol.Transaction{protocol.NewTransaction(recipientAddress, senderAddress, nil, 0, 1)}
	block := protocol.NewBlock(0, [32]byte{}, transactions, nil)
	isRegistered, _ := block.IsValidatorRegistered()
	test.Assert(t, isRegistered, "proof of humanity is invalid whereas it should be")
}
