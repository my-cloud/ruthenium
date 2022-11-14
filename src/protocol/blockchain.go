package protocol

import (
	"github.com/my-cloud/ruthenium/src/network"
)

type Blockchain interface {
	Verify()
	StartVerification()
	Blocks() []*network.BlockResponse
	CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64
	AddBlock(blockResponse *network.BlockResponse)
}
