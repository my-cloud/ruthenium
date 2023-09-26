package protocoltest

//
//func NewGenesisBlockResponse(validatorWalletAddress string, genesisValue uint64) *network.BlockResponse {
//	genesisTransaction, _ := validation.NewRewardTransaction(validatorWalletAddress, true, 0, genesisValue)
//	return &network.BlockResponse{
//		Timestamp:                  0,
//		PreviousHash:               [32]byte{},
//		Transactions:               []*validation.Transaction{genesisTransaction},
//		AddedRegisteredAddresses:   nil,
//		RemovedRegisteredAddresses: nil,
//	}
//}
//
//func NewRewardedBlockResponse(previousHash [32]byte, timestamp int64) *network.BlockResponse {
//	rewardTransaction, _ := validation.NewRewardTransaction("recipient", false, 0, 0)
//	return &network.BlockResponse{
//		Timestamp:                  timestamp,
//		PreviousHash:               previousHash,
//		Transactions:               []*network.TransactionResponse{rewardTransaction},
//		AddedRegisteredAddresses:   nil,
//		RemovedRegisteredAddresses: nil,
//	}
//}
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
