package network

type UtxosRequest struct {
	Address *string
}

func (utxosRequest UtxosRequest) IsInvalid() bool {
	return utxosRequest.Address == nil || len(*utxosRequest.Address) == 0
}
