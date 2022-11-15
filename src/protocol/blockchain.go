package protocol

import (
	"github.com/my-cloud/ruthenium/src/network"
)

type Blockchain interface {
	Verify(timestamp int64)
	Blocks() []*network.BlockResponse
	CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64
	AddBlock(timestamp int64, previousHash [32]byte, transactions []*network.TransactionResponse, registeredAddresses []string)
	Copy() Blockchain
	IsEmpty() bool
}
