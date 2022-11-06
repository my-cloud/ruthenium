package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/protocol"
)

func NewGenesisBlockResponse(validatorWalletAddress string) *neighborhood.BlockResponse {
	genesisTransaction := protocol.NewRewardTransaction(validatorWalletAddress, 0, 1e13)
	return &neighborhood.BlockResponse{
		Timestamp:           0,
		PreviousHash:        [32]byte{},
		Transactions:        []*neighborhood.TransactionResponse{genesisTransaction.GetResponse()},
		RegisteredAddresses: nil,
	}
}

func NewRewardedBlockResponse(previousHash [32]byte, timestamp int64) *neighborhood.BlockResponse {
	rewardTransaction := protocol.NewRewardTransaction("recipient", 0, 0)
	return &neighborhood.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        previousHash,
		Transactions:        []*neighborhood.TransactionResponse{rewardTransaction.GetResponse()},
		RegisteredAddresses: nil,
	}
}

func NewEmptyBlockResponse(timestamp int64) *neighborhood.BlockResponse {
	return &neighborhood.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        [32]byte{},
		Transactions:        nil,
		RegisteredAddresses: nil,
	}
}

func NewBlockResponse(timestamp int64, transaction *protocol.Transaction) *neighborhood.BlockResponse {
	return &neighborhood.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        [32]byte{},
		Transactions:        []*neighborhood.TransactionResponse{transaction.GetResponse()},
		RegisteredAddresses: nil,
	}
}
