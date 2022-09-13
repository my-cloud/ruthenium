package blockchain

import (
	"gitlab.com/coinsmaster/ruthenium/src/node/blockchain"
	"gitlab.com/coinsmaster/ruthenium/test"
	"testing"
)

func Test_IsProofOfHumanityValid_Invalid_ReturnsFalse(t *testing.T) {
	senderAddress := blockchain.RewardSenderAddress
	recipientAddress := "0x1234567890123456789012345678901234567890"
	transactions := []*blockchain.Transaction{blockchain.NewTransaction(recipientAddress, senderAddress, nil, 0, 1)}
	block := blockchain.NewBlock(0, [32]byte{}, transactions)
	err := block.IsProofOfHumanityValid()
	test.Assert(t, err != nil, "proof of humanity is valid whereas it should not")
}

func Test_IsProofOfHumanityValid_Valid_ReturnsTrue(t *testing.T) {
	senderAddress := blockchain.RewardSenderAddress
	recipientAddress := "0xD70eFbfC9cF73dC1aB5f8292E7273D916D38c916"
	transactions := []*blockchain.Transaction{blockchain.NewTransaction(recipientAddress, senderAddress, nil, 0, 1)}
	block := blockchain.NewBlock(0, [32]byte{}, transactions)
	err := block.IsProofOfHumanityValid()
	test.Assert(t, err == nil, "proof of humanity is invalid whereas it should be")
}
