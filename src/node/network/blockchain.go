package network

import (
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
)

type Blockchain interface {
	Verify()
	StartVerification()
	Blocks() []*neighborhood.BlockResponse
	CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64
	AddBlock(blockResponse *neighborhood.BlockResponse)
}
