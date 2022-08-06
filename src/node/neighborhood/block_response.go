package neighborhood

import "ruthenium/src/node/authentication"

type BlockResponse struct {
	Timestamp    int64
	Nonce        int
	PreviousHash [32]byte
	Transactions []*authentication.TransactionResponse
}
