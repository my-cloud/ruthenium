package network

type BlocksRequest struct {
	StartingBlockHeight *uint64
}

func (blocksRequest BlocksRequest) IsInvalid() bool {
	return blocksRequest.StartingBlockHeight == nil
}
