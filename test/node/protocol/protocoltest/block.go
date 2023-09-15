package protocoltest

import (
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
)

func NewGenesisBlockResponse(validatorWalletAddress string, genesisValue uint64) *network.BlockResponse {
	genesisTransaction, _ := validation.NewGenesisTransaction(validatorWalletAddress, 0, genesisValue)
	return &network.BlockResponse{
		Timestamp:                  0,
		PreviousHash:               [32]byte{},
		Transactions:               []*network.TransactionResponse{genesisTransaction},
		AddedRegisteredAddresses:   nil,
		RemovedRegisteredAddresses: nil,
	}
}

func NewRewardedBlockResponse(previousHash [32]byte, timestamp int64) *network.BlockResponse {
	rewardTransaction, _ := validation.NewRewardTransaction("recipient", 0, 0)
	return &network.BlockResponse{
		Timestamp:                  timestamp,
		PreviousHash:               previousHash,
		Transactions:               []*network.TransactionResponse{rewardTransaction},
		AddedRegisteredAddresses:   nil,
		RemovedRegisteredAddresses: nil,
	}
}

func NewEmptyBlockResponse(timestamp int64) *network.BlockResponse {
	return &network.BlockResponse{
		Timestamp:                  timestamp,
		PreviousHash:               [32]byte{},
		Transactions:               nil,
		AddedRegisteredAddresses:   nil,
		RemovedRegisteredAddresses: nil,
	}
}
