package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Blockchain interface {
	AddBlock(timestamp int64, transactions []*network.TransactionResponse, registeredAddresses []string)
	Blocks() []*network.BlockResponse
	CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64
	Copy() Blockchain
	LastBlocks(startingBlockHash *[32]byte) []*network.BlockResponse
	Update(timestamp int64)
}
