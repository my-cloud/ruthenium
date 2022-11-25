package protocol

import (
	network2 "github.com/my-cloud/ruthenium/src/node/network"
)

type Blockchain interface {
	Verify(timestamp int64)
	Blocks() []*network2.BlockResponse
	CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64
	AddBlock(timestamp int64, transactions []*network2.TransactionResponse, registeredAddresses []string)
	Copy() Blockchain
	IsEmpty() bool
}
