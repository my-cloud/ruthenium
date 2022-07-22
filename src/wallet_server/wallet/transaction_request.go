package wallet

type TransactionRequest struct {
	SenderPrivateKey *string `json:"sender_private_key"`
	SenderAddress    *string `json:"sender_address"`
	RecipientAddress *string `json:"recipient_address"`
	SenderPublicKey  *string `json:"sender_public_key"`
	Value            *string `json:"value"`
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.SenderPrivateKey == nil ||
		transactionRequest.SenderAddress == nil ||
		transactionRequest.RecipientAddress == nil ||
		transactionRequest.SenderPublicKey == nil ||
		transactionRequest.Value == nil
}
