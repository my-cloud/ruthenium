package network

type UtxoResponse struct {
	Address       string
	BlockHeight   int // TODO uint64
	HasReward     bool
	HasIncome     bool
	OutputIndex   uint16
	TransactionId string
	Value         uint64
}
