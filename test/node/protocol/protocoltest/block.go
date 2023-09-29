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

//
//func NewEmptyBlockResponse(timestamp int64) *network.BlockResponse {
//	return &network.BlockResponse{
//		Timestamp:                  timestamp,
//		PreviousHash:               [32]byte{},
//		Transactions:               nil,
//		AddedRegisteredAddresses:   nil,
//		RemovedRegisteredAddresses: nil,
//	}
//}
