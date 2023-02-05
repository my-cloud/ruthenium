package network

type TransactionRequest struct {
	Fee                          *uint64
	RecipientAddress             *string
	SenderAddress                *string
	SenderPublicKey              *string
	Signature                    *string
	Timestamp                    *int64
	TransactionBroadcasterTarget *string
	Value                        *uint64
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.Fee == nil ||
		transactionRequest.RecipientAddress == nil || len(*transactionRequest.RecipientAddress) == 0 ||
		transactionRequest.SenderAddress == nil || len(*transactionRequest.SenderAddress) == 0 ||
		transactionRequest.SenderPublicKey == nil || len(*transactionRequest.SenderPublicKey) == 0 ||
		transactionRequest.Signature == nil || len(*transactionRequest.Signature) == 0 ||
		transactionRequest.Timestamp == nil ||
		transactionRequest.TransactionBroadcasterTarget == nil || len(*transactionRequest.TransactionBroadcasterTarget) == 0 ||
		transactionRequest.Value == nil
}
