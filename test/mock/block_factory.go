package mock

import (
	network2 "github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
)

func NewGenesisBlockResponse(validatorWalletAddress string) *network2.BlockResponse {
	genesisTransaction := validation.NewRewardTransaction(validatorWalletAddress, 0, 1e13)
	return &network2.BlockResponse{
		Timestamp:           0,
		PreviousHash:        [32]byte{},
		Transactions:        []*network2.TransactionResponse{genesisTransaction},
		RegisteredAddresses: nil,
	}
}

func NewRewardedBlockResponse(previousHash [32]byte, timestamp int64) *network2.BlockResponse {
	rewardTransaction := validation.NewRewardTransaction("recipient", 0, 0)
	return &network2.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        previousHash,
		Transactions:        []*network2.TransactionResponse{rewardTransaction},
		RegisteredAddresses: nil,
	}
}

func NewEmptyBlockResponse(timestamp int64) *network2.BlockResponse {
	return &network2.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        [32]byte{},
		Transactions:        nil,
		RegisteredAddresses: nil,
	}
}

func NewBlockResponse(timestamp int64, hash [32]byte, transactionResponses []*network2.TransactionResponse, registeredAddresses []string) *network2.BlockResponse {
	return &network2.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        hash,
		Transactions:        transactionResponses,
		RegisteredAddresses: registeredAddresses,
	}
}
