package network

type TransactionRequest struct {
	Inputs                       *[]InputRequest  `json:"inputs"`
	Outputs                      *[]OutputRequest `json:"outputs"`
	Timestamp                    *int64           `json:"timestamp"`
	TransactionBroadcasterTarget *string
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.Inputs == nil || len(*transactionRequest.Inputs) == 0 ||
		transactionRequest.Outputs == nil || len(*transactionRequest.Outputs) == 0 ||
		transactionRequest.Timestamp == nil ||
		transactionRequest.TransactionBroadcasterTarget == nil || len(*transactionRequest.TransactionBroadcasterTarget) == 0
}
