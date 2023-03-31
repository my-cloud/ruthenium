package network

type TransactionRequest struct {
	Inputs                       *[]InputRequest  `json:"inputs"`
	Outputs                      *[]OutputRequest `json:"outputs"`
	Timestamp                    *int64           `json:"timestamp"`
	TransactionBroadcasterTarget *string
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	if transactionRequest.Inputs == nil {
		return true
	}
	for _, input := range *transactionRequest.Inputs {
		if input.IsInvalid() {
			return true
		}
	}
	if transactionRequest.Outputs == nil || len(*transactionRequest.Outputs) == 0 {
		return true
	}
	for _, output := range *transactionRequest.Outputs {
		if output.IsInvalid() {
			return true
		}
	}
	return transactionRequest.Timestamp == nil ||
		transactionRequest.TransactionBroadcasterTarget == nil || len(*transactionRequest.TransactionBroadcasterTarget) == 0
}
