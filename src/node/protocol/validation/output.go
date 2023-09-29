package validation

import (
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/node/network"
)

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
	return json.Marshal(network.OutputResponse{
		Address:   output.address,
		HasIncome: output.hasIncome,
		HasReward: output.hasReward,
		Value:     output.value,
	})
}

func (output *Output) UnmarshalJSON(data []byte) error {
	var outputDto network.OutputResponse
	err := json.Unmarshal(data, &outputDto)
	if err != nil {
		return err
	}
	output.address = outputDto.Address
	output.hasIncome = outputDto.HasIncome
	output.hasReward = outputDto.HasReward
	output.value = outputDto.Value
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
