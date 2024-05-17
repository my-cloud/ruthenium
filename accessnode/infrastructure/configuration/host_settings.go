package configuration

import (
	"encoding/json"
	"strconv"
)

type hostSettingsDto struct {
	Port int
}

type HostSettings struct {
	port string
}

func (settings *HostSettings) UnmarshalJSON(data []byte) error {
	var dto *hostSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.port = strconv.Itoa(dto.Port)
	return nil
}

func (settings *HostSettings) Port() string {
	return settings.port
}
