package server

type TransactionRequest struct {
	RecipientAddress *string `json:"recipient_address"`
	Value            *string `json:"value"`
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.RecipientAddress == nil || len(*transactionRequest.RecipientAddress) == 0 ||
		transactionRequest.Value == nil || len(*transactionRequest.Value) == 0
}
