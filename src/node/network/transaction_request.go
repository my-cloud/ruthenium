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
	for _, i := range *transactionRequest.Inputs {
		input := &i
		if input == nil || input.IsInvalid() {
			return true
		}
	}
	if transactionRequest.Outputs == nil || len(*transactionRequest.Outputs) == 0 {
		return true
	}
	for _, o := range *transactionRequest.Outputs {
		output := &o
		if output == nil || output.IsInvalid() {
			return true
		}
	}
	return transactionRequest.Timestamp == nil ||
		transactionRequest.TransactionBroadcasterTarget == nil || len(*transactionRequest.TransactionBroadcasterTarget) == 0
}
