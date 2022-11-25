package network

type AmountRequest struct {
	Address *string
}

func (amountRequest *AmountRequest) IsInvalid() bool {
	return amountRequest.Address == nil || len(*amountRequest.Address) == 0
}
