package network

type TargetRequest struct {
	Target *string
}

func (targetRequest TargetRequest) IsInvalid() bool {
	return targetRequest.Target == nil || len(*targetRequest.Target) == 0
}
