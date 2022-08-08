package neighborhood

type TransactionResponse struct {
	Timestamp        int64
	SenderAddress    string
	RecipientAddress string
	Value            float32
}
