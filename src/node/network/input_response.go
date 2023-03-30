package network

type InputResponse struct {
	OutputIndex   uint16
	TransactionId [32]byte
	PublicKey     string
	Signature     string
}
