package neighborhood

type TargetRequest struct {
	Ip   *string
	Port *uint16
}

func (targetRequest *TargetRequest) IsInvalid() bool {
	return targetRequest.Ip == nil || len(*targetRequest.Ip) == 0 ||
		targetRequest.Port == nil
}
