package network

type TransactionRequest struct {
	RecipientAddress             *string `json:"recipient_address"`
	SenderAddress                *string `json:"sender_address"`
	SenderPublicKey              *string `json:"sender_public_key"`
	Signature                    *string `json:"signature"`
	Timestamp                    *int64  `json:"timestamp"`
	TransactionBroadcasterTarget *string
	Value                        *uint64 `json:"value"`
	Fee                          *uint64 `json:"fee"`
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.RecipientAddress == nil || len(*transactionRequest.RecipientAddress) == 0 ||
		transactionRequest.SenderAddress == nil || len(*transactionRequest.SenderAddress) == 0 ||
		transactionRequest.SenderPublicKey == nil || len(*transactionRequest.SenderPublicKey) == 0 ||
		transactionRequest.Signature == nil || len(*transactionRequest.Signature) == 0 ||
		transactionRequest.Timestamp == nil ||
		transactionRequest.TransactionBroadcasterTarget == nil || len(*transactionRequest.TransactionBroadcasterTarget) == 0 ||
		transactionRequest.Value == nil ||
		transactionRequest.Fee == nil
}
