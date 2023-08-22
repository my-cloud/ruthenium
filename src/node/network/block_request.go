package network

type BlocksRequest struct {
	BlockHeight *uint64
}

func (blocksRequest BlocksRequest) IsInvalid() bool {
	return blocksRequest.BlockHeight == nil
}
