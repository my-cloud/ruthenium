package protocol

import (
	"github.com/my-cloud/ruthenium/src/api/node"
)

type Verifiable interface {
	Verify()
	StartVerification(validatable Validatable)
	Blocks() []*node.BlockResponse
	CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64
	AddBlock(blockResponse *node.BlockResponse)
}
