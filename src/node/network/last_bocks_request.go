package network

type LastBlocksRequest struct {
	StartingBlockNonce *int
}

func (lastBlocksRequest LastBlocksRequest) IsInvalid() bool {
	return lastBlocksRequest.StartingBlockNonce == nil
}
