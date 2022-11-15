package verification

import (
	"github.com/my-cloud/ruthenium/src/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
)

func NewGenesisBlockResponse(validatorWalletAddress string) *network.BlockResponse {
	genesisTransaction := validation.NewRewardTransaction(validatorWalletAddress, 0, 1e13)
	return &network.BlockResponse{
		Timestamp:           0,
		PreviousHash:        [32]byte{},
		Transactions:        []*network.TransactionResponse{genesisTransaction.GetResponse()},
		RegisteredAddresses: nil,
	}
}

func NewRewardedBlockResponse(previousHash [32]byte, timestamp int64) *network.BlockResponse {
	rewardTransaction := validation.NewRewardTransaction("recipient", 0, 0)
	return &network.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        previousHash,
		Transactions:        []*network.TransactionResponse{rewardTransaction.GetResponse()},
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

func NewBlockResponse(timestamp int64, hash [32]byte, transactions ...*validation.Transaction) *network.BlockResponse {
	var transactionResponses []*network.TransactionResponse
	var registeredAddresses []string
	registeredAddressesMap := make(map[string]bool)
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, transaction.GetResponse())
		if _, ok := registeredAddressesMap[transaction.SenderAddress()]; !ok && !transaction.IsReward() {
			registeredAddressesMap[transaction.SenderAddress()] = true
		}
	}
	for address := range registeredAddressesMap {
		registeredAddresses = append(registeredAddresses, address)
	}
	return &network.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        hash,
		Transactions:        transactionResponses,
		RegisteredAddresses: registeredAddresses,
	}
}
