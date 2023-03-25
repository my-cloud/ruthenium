package network

type InputResponse struct {
	OutputIndex           uint16
	PreviousTransactionId [32]byte
	PublicKey             string
	Signature             string
}
