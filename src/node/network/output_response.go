package network

type OutputResponse struct {
	Address   string `json:"address"`
	HasReward bool   `json:"has_reward"`
	HasIncome bool   `json:"has_income"`
	Value     uint64 `json:"value"`
}
