package configuration

import (
	"encoding/json"
	"strconv"
)

type hostSettingsDto struct {
	Ip   string
	Port int
}

type HostSettings struct {
	ip   string
	port string
}

func (settings *HostSettings) UnmarshalJSON(data []byte) error {
	var dto *hostSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.ip = dto.Ip
	settings.port = strconv.Itoa(dto.Port)
	return nil
}

func (settings *HostSettings) Ip() string {
	return settings.ip
}

func (settings *HostSettings) Port() string {
	return settings.port
}
