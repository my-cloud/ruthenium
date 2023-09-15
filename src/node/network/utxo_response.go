package network

type UtxoResponse struct {
	Address       string
	BlockHeight   int
	HasReward     bool
	HasIncome     bool
	OutputIndex   uint16
	TransactionId string
	Value         uint64
}
