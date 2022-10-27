package neighborhood

type TransactionRequest struct {
	RecipientAddress *string
	SenderAddress    *string
	SenderPublicKey  *string
	Signature        *string
	Timestamp        *int64
	Value            *uint64
	Fee              *uint64
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.RecipientAddress == nil || len(*transactionRequest.RecipientAddress) == 0 ||
		transactionRequest.SenderAddress == nil || len(*transactionRequest.SenderAddress) == 0 ||
		transactionRequest.SenderPublicKey == nil || len(*transactionRequest.SenderPublicKey) == 0 ||
		transactionRequest.Signature == nil || len(*transactionRequest.Signature) == 0 ||
		transactionRequest.Timestamp == nil ||
		transactionRequest.Value == nil ||
		transactionRequest.Fee == nil
}
