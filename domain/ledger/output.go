package ledger

import (
	"encoding/json"
)

type outputDto struct {
	Address    string `json:"address"`
	IsYielding bool   `json:"is_yielding"`
	Value      uint64 `json:"value"`
}

type Output struct {
	address    string
	isYielding bool
	value      uint64
}

func NewOutput(address string, isYielding bool, value uint64) *Output {
	return &Output{address, isYielding, value}
}

func (output *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(outputDto{
		Address:    output.address,
		IsYielding: output.isYielding,
		Value:      output.value,
	})
}

func (output *Output) UnmarshalJSON(data []byte) error {
	var dto *outputDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	output.address = dto.Address
	output.isYielding = dto.IsYielding
	output.value = dto.Value
	return nil
}

func (output *Output) Address() string {
	return output.address
}

func (output *Output) IsYielding() bool {
	return output.isYielding
}

func (output *Output) InitialValue() uint64 {
	return output.value
}
