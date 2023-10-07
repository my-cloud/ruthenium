package protocoltest

import (
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
)

func NewGenesisBlock(validatorWalletAddress string, genesisValue uint64) *verification.Block {
	genesisTransaction, _ := verification.NewRewardTransaction(validatorWalletAddress, true, 0, genesisValue)
	transactions := []*verification.Transaction{genesisTransaction}
	return verification.NewBlock(0, [32]byte{}, transactions, nil, nil)
}

func NewRewardedBlock(previousHash [32]byte, timestamp int64) *verification.Block {
	rewardTransaction, _ := verification.NewRewardTransaction("recipient", false, 0, 0)
	transactions := []*verification.Transaction{rewardTransaction}
	return verification.NewBlock(timestamp, previousHash, transactions, nil, nil)
}
