package chain

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type TransactionRequest struct {
	Verb             *string
	SenderAddress    *string
	RecipientAddress *string
	SenderPublicKey  *string
	Value            *float32
	Signature        *string
}

func (transactionRequest *TransactionRequest) IsInvalid() bool {
	return transactionRequest.Verb == nil ||
		transactionRequest.SenderAddress == nil ||
		transactionRequest.RecipientAddress == nil ||
		transactionRequest.SenderPublicKey == nil ||
		transactionRequest.Value == nil
}
