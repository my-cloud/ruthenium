package network

type UtxoResponse struct {
	Address       string `json:"address"`
	BlockHeight   int    `json:"block_height"`
	HasReward     bool   `json:"has_reward"`
	HasIncome     bool   `json:"has_income"`
	OutputIndex   uint16 `json:"output_index"`
	TransactionId string `json:"transaction_id"`
	Value         uint64 `json:"value"`
}
