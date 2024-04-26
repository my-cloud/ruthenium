package configuration

import (
	"encoding/json"
)

type validatorSettingsDto struct {
	Address   string
	InfuraKey string
}

type ValidatorSettings struct {
	address   string
	infuraKey string
}

func (settings *ValidatorSettings) UnmarshalJSON(data []byte) error {
	var dto *validatorSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.address = dto.Address
	settings.infuraKey = dto.InfuraKey
	return nil
}

func (settings *ValidatorSettings) Address() string {
	return settings.address
}

func (settings *ValidatorSettings) InfuraKey() string {
	return settings.infuraKey
}
