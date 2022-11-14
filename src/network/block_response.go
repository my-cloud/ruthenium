package network

type BlockResponse struct {
	Timestamp           int64
	PreviousHash        [32]byte
	Transactions        []*TransactionResponse
	RegisteredAddresses []string
}
