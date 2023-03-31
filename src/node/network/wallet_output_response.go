package network

type WalletOutputResponse struct {
	Address       string   `json:"address"`
	BlockHeight   int      `json:"block_height"`
	HasReward     bool     `json:"has_reward"`
	HasIncome     bool     `json:"has_income"`
	OutputIndex   uint16   `json:"output_index"`
	TransactionId [32]byte `json:"transaction_id"`
	Value         uint64   `json:"value"`
}
