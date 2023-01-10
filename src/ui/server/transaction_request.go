package server

type TransactionRequest struct {
	SenderPrivateKey *string `json:"sender_private_key"`
	SenderAddress    *string `json:"sender_address"`
	RecipientAddress *string `json:"recipient_address"`
	Value            *string `json:"value"`
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.SenderPrivateKey == nil || len(*transactionRequest.SenderPrivateKey) == 0 ||
		transactionRequest.SenderAddress == nil || len(*transactionRequest.SenderAddress) == 0 ||
		transactionRequest.RecipientAddress == nil || len(*transactionRequest.RecipientAddress) == 0 ||
		transactionRequest.Value == nil || len(*transactionRequest.Value) == 0
}
