package protocoltest

import (
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
)

func NewGenesisBlock(validatorWalletAddress string, genesisValue uint64) *verification.Block {
	genesisTransaction, _ := validation.NewRewardTransaction(validatorWalletAddress, true, 0, genesisValue)
	transactions := []*validation.Transaction{genesisTransaction}
	return verification.NewBlock(0, [32]byte{}, transactions, nil, nil)
}

func NewRewardedBlock(previousHash [32]byte, timestamp int64) *verification.Block {
	rewardTransaction, _ := validation.NewRewardTransaction("recipient", false, 0, 0)
	transactions := []*validation.Transaction{rewardTransaction}
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
