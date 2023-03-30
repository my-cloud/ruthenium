package protocol

import (
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Blockchain interface {
	AddBlock(timestamp int64, transactions []*network.TransactionResponse, newRegisteredAddresses []string) error
	Block(blockHeight uint64) *network.BlockResponse
	Blocks() []*network.BlockResponse
	CalculateFee(transaction *network.TransactionResponse, timestamp int64) uint64
	Copy() Blockchain
	Lambda() float64
	LastBlocks(startingBlockHeight uint64) []*network.BlockResponse
	UtxosByAddress(address string) []*network.OutputResponse
}
