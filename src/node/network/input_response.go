package network

type InputResponse struct {
	OutputIndex   uint16   `json:"output-index"`
	TransactionId [32]byte `json:"transaction_id"`
	PublicKey     string   `json:"public_key"`
	Signature     string   `json:"signature"`
}
