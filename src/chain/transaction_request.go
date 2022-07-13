package chain

type TransactionRequest struct {
	SenderAddress    *string  `json:"sender_address"`
	RecipientAddress *string  `json:"recipient_address"`
	SenderPublicKey  *string  `json:"sender_public_key"`
	Value            *float32 `json:"value"`
	Signature        *string  `json:"signature"`
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.SenderAddress == nil ||
		transactionRequest.RecipientAddress == nil ||
		transactionRequest.SenderPublicKey == nil ||
		transactionRequest.Value == nil
}
