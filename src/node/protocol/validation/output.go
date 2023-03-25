package validation

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
)

type Output struct {
	address     string
	blockHeight int
	hasReward   bool
	hasIncome   bool
	value       uint64
}

func NewOutput(address string, blockHeight int, hasReward bool, hasIncome bool, value uint64) *Output {
	return &Output{address, blockHeight, hasReward, hasIncome, value}
}

func (output *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address     string `json:"address"`
		BlockHeight bool   `json:"block_height"`
		HasReward   bool   `json:"has_reward"`
		HasIncome   bool   `json:"has_income"`
		Value       uint64 `json:"value"`
	}{
		Address:   output.address,
		HasReward: output.hasReward,
		HasIncome: output.hasIncome,
		Value:     output.value,
	})
}

func (output *Output) GetResponse() *network.OutputResponse {
	return &network.OutputResponse{
		Address:     output.address,
		BlockHeight: output.blockHeight,
		HasReward:   output.hasReward,
		HasIncome:   output.hasIncome,
		Value:       output.value,
	}
}
