package network

type OutputRequest struct {
	Address   *string `json:"address"`
	HasReward *bool   `json:"has_reward"`
	HasIncome *bool   `json:"has_income"`
	Value     *uint64 `json:"value"`
}

func (outputRequest *OutputRequest) IsInvalid() bool {
	return outputRequest.Address == nil || len(*outputRequest.Address) == 0 ||
		outputRequest.HasReward == nil ||
		outputRequest.HasIncome == nil ||
		outputRequest.Value == nil
}
