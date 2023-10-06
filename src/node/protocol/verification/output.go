package verification

import (
	"encoding/json"
)

type outputDto struct {
	Address      string `json:"address"`
	IsRegistered bool   `json:"is_registered"`
	Value        uint64 `json:"value"`
}

type Output struct {
	address      string
	isRegistered bool
	value        uint64
}

func NewOutput(address string, isRegistered bool, value uint64) *Output {
	return &Output{address, isRegistered, value}
}

func (output *Output) MarshalJSON() ([]byte, error) {
	return json.Marshal(outputDto{
		Address:      output.address,
		IsRegistered: output.isRegistered,
		Value:        output.value,
	})
}

func (output *Output) UnmarshalJSON(data []byte) error {
	var dto *outputDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	output.address = dto.Address
	output.isRegistered = dto.IsRegistered
	output.value = dto.Value
	return nil
}

func (output *Output) Address() string {
	return output.address
}

func (output *Output) IsRegistered() bool {
	return output.isRegistered
}

func (output *Output) InitialValue() uint64 {
	return output.value
}
