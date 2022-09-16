package blockchain

import (
	"gitlab.com/coinsmaster/ruthenium/src/node/blockchain"
	"gitlab.com/coinsmaster/ruthenium/test"
	"testing"
)

func Test_IsProofOfHumanityValid_Invalid_ReturnsFalse(t *testing.T) {
	senderAddress := blockchain.RewardSenderAddress
	recipientAddress := "0x0000000000000000000000000000000000000001"
	transactions := []*blockchain.Transaction{blockchain.NewTransaction(recipientAddress, senderAddress, nil, 0, 1)}
	block := blockchain.NewBlock(0, [32]byte{}, transactions)
	err := block.IsProofOfHumanityValid()
	test.Assert(t, err != nil, "proof of humanity is valid whereas it should not")
}

func Test_IsProofOfHumanityValid_Valid_ReturnsTrue(t *testing.T) {
	senderAddress := blockchain.RewardSenderAddress
	recipientAddress := "0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a"
	transactions := []*blockchain.Transaction{blockchain.NewTransaction(recipientAddress, senderAddress, nil, 0, 1)}
	block := blockchain.NewBlock(0, [32]byte{}, transactions)
	err := block.IsProofOfHumanityValid()
	test.Assert(t, err == nil, "proof of humanity is invalid whereas it should be")
}
