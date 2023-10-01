package verification

import (
	"encoding/json"
)

type outputDto struct {
	Address   string `json:"address"`
	HasIncome bool   `json:"has_income"`
	HasReward bool   `json:"has_reward"`
	Value     uint64 `json:"value"`
}

type Output struct {
	address   string
	hasIncome bool
	hasReward bool
	value     uint64
}

func NewOutput(address string, hasIncome bool, hasReward bool, value uint64) *Output {
	return &Output{address, hasIncome, hasReward, value}
}

func (output *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(outputDto{
		Address:   output.address,
		HasIncome: output.hasIncome,
		HasReward: output.hasReward,
		Value:     output.value,
	})
}

func (output *Output) UnmarshalJSON(data []byte) error {
	var dto *outputDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	output.address = dto.Address
	output.hasIncome = dto.HasIncome
	output.hasReward = dto.HasReward
	output.value = dto.Value
	return nil
}

func (output *Output) Address() string {
	return output.address
}

func (output *Output) HasIncome() bool {
	return output.hasIncome
}

func (output *Output) HasReward() bool {
	return output.hasReward
}

func (output *Output) InitialValue() uint64 {
	return output.value
}
