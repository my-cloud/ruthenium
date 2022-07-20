package chain

type PostTransactionRequest struct {
	SenderAddress    *string
	RecipientAddress *string
	SenderPublicKey  *string
	Value            *float32
	Signature        *string
}

func (transactionRequest *PostTransactionRequest) IsInvalid() bool {
	return transactionRequest.SenderAddress == nil ||
		transactionRequest.RecipientAddress == nil ||
		transactionRequest.SenderPublicKey == nil ||
		transactionRequest.Value == nil
}
