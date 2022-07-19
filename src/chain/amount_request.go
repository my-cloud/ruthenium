package chain

type AmountRequest struct {
	Address *string `json:"address"`
}

func (amountRequest *AmountRequest) IsInvalid() bool {
	return amountRequest.Address == nil
}
