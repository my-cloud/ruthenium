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

func NewBlockResponse(timestamp int64, hash [32]byte, transactions ...*protocol.Transaction) *neighborhood.BlockResponse {
	var transactionResponses []*neighborhood.TransactionResponse
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
	return &neighborhood.BlockResponse{
		Timestamp:           timestamp,
		PreviousHash:        hash,
		Transactions:        transactionResponses,
		RegisteredAddresses: registeredAddresses,
	}
}
