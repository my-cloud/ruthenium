package protocol

func NewGenesisBlock(validatorWalletAddress string, genesisValue uint64) *Block {
	genesisTransaction, _ := NewRewardTransaction(validatorWalletAddress, true, 0, genesisValue)
	transactions := []*Transaction{genesisTransaction}
	return NewBlock([32]byte{}, nil, nil, 0, transactions)
}

func NewRewardedBlock(previousHash [32]byte, timestamp int64) *Block {
	rewardTransaction, _ := NewRewardTransaction("recipient", false, 0, 0)
	transactions := []*Transaction{rewardTransaction}
	return NewBlock(previousHash, nil, nil, timestamp, transactions)
}
