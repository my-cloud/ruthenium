package network

type InputRequest struct {
	OutputIndex   *uint16 `json:"output_index"`
	TransactionId *string `json:"transaction_id"`
	PublicKey     *string `json:"public_key"`
	Signature     *string `json:"signature"`
}

func (inputRequest *InputRequest) IsInvalid() bool {
	return inputRequest.OutputIndex == nil ||
		inputRequest.TransactionId == nil || len(*inputRequest.TransactionId) == 0 ||
		inputRequest.PublicKey == nil || len(*inputRequest.PublicKey) == 0 ||
		inputRequest.Signature == nil || len(*inputRequest.Signature) == 0
}
