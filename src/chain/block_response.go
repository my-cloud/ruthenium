package chain

type BlockResponse struct {
	Timestamp    int64
	Nonce        int
	PreviousHash [32]byte
	Transactions []*TransactionResponse
}
