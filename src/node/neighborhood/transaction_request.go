package neighborhood

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type TransactionRequest struct {
	RecipientAddress *string
	SenderAddress    *string
	SenderPublicKey  *string
	Signature        *string
	Timestamp        *int64
	Value            *uint64
	Verb             *string
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.RecipientAddress == nil || len(*transactionRequest.RecipientAddress) == 0 ||
		transactionRequest.SenderAddress == nil || len(*transactionRequest.SenderAddress) == 0 ||
		transactionRequest.SenderPublicKey == nil || len(*transactionRequest.SenderPublicKey) == 0 ||
		transactionRequest.Signature == nil || len(*transactionRequest.Signature) == 0 ||
		transactionRequest.Timestamp == nil ||
		transactionRequest.Value == nil ||
		transactionRequest.Verb == nil || len(*transactionRequest.Verb) == 0
}
