package authentication

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
	return transactionRequest.Verb == nil || len(*transactionRequest.Verb) == 0 ||
		transactionRequest.SenderAddress == nil || len(*transactionRequest.SenderAddress) == 0 ||
		transactionRequest.RecipientAddress == nil || len(*transactionRequest.RecipientAddress) == 0 ||
		transactionRequest.SenderPublicKey == nil || len(*transactionRequest.SenderPublicKey) == 0 ||
		transactionRequest.Value == nil ||
		transactionRequest.Signature == nil || len(*transactionRequest.Signature) == 0
}
