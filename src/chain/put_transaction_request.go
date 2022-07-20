package chain

type PutTransactionRequest struct {
	SenderAddress    *string
	RecipientAddress *string
	SenderPublicKey  *string
	Value            *float32
	Signature        *string
}

func (transactionRequest *PutTransactionRequest) IsInvalid() bool {
	return transactionRequest.SenderAddress == nil ||
		transactionRequest.RecipientAddress == nil ||
		transactionRequest.SenderPublicKey == nil ||
		transactionRequest.Value == nil
}
