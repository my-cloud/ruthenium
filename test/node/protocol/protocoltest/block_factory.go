package protocoltest

import (
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
)

func NewGenesisBlockResponse(validatorWalletAddress string) *network.BlockResponse {
	genesisTransaction := validation.NewRewardTransaction(validatorWalletAddress, 0, 1e13)
	return &network.BlockResponse{
		Timestamp:           0,
		PreviousHash:        [32]byte{},
		Transactions:        []*network.TransactionResponse{genesisTransaction},
		RegisteredAddresses: nil,
	}
}

func NewRewardedBlockResponse(previousHash [32]byte, timestamp int64) *network.BlockResponse {
	rewardTransaction := validation.NewRewardTransaction("recipient", 0, 0)
	return &network.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        previousHash,
		Transactions:        []*network.TransactionResponse{rewardTransaction},
		RegisteredAddresses: nil,
	}
}

func NewEmptyBlockResponse(timestamp int64) *network.BlockResponse {
	return &network.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        [32]byte{},
		Transactions:        nil,
		RegisteredAddresses: nil,
	}
}

func NewBlockResponse(timestamp int64, hash [32]byte, transactionResponses []*network.TransactionResponse, registeredAddresses []string) *network.BlockResponse {
	return &network.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        hash,
		Transactions:        transactionResponses,
		RegisteredAddresses: registeredAddresses,
	}
}
