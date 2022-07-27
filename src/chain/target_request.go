package chain

type TargetRequest struct {
	Kind *string
	Ip   *string
	Port *uint16
}

func (targetRequest *TargetRequest) IsInvalid() bool {
	return targetRequest.Kind == nil || len(*targetRequest.Kind) == 0 ||
		targetRequest.Ip == nil || len(*targetRequest.Ip) == 0 ||
		targetRequest.Port == nil
}
