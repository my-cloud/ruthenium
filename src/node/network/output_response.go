package network

type OutputResponse struct {
	Address   string `json:"address"`
	HasIncome bool   `json:"has_income"`
	HasReward bool   `json:"has_reward"`
	Value     uint64 `json:"value"`
}
