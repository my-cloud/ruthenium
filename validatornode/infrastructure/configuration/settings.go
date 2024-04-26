package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type settingsDto struct {
	Host      *HostSettings
	Network   *NetworkSettings
	Protocol  *ProtocolSettings
	Registry  *RegistrySettings
	Validator *ValidatorSettings
	Log       *LogSettings
}

type Settings struct {
	host      *HostSettings
	network   *NetworkSettings
	protocol  *ProtocolSettings
	registry  *RegistrySettings
	validator *ValidatorSettings
	log       *LogSettings
}

func NewSettings(path string) (*Settings, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	var settings *Settings
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}
	if err = jsonFile.Close(); err != nil {
		return nil, fmt.Errorf("unable to close file: %w", err)
	}
	if err = json.Unmarshal(bytes, &settings); err != nil {
		return nil, fmt.Errorf("unable to unmarshal: %w", err)
	}
	return settings, nil
}

func (settings *Settings) UnmarshalJSON(data []byte) error {
	var dto *settingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.host = dto.Host
	settings.network = dto.Network
	settings.protocol = dto.Protocol
	settings.registry = dto.Registry
	settings.validator = dto.Validator
	settings.log = dto.Log
	return nil
}

func (settings *Settings) Host() *HostSettings {
	return settings.host
}

func (settings *Settings) Network() *NetworkSettings {
	return settings.network
}

func (settings *Settings) Protocol() *ProtocolSettings {
	return settings.protocol
}

func (settings *Settings) Registry() *RegistrySettings {
	return settings.registry
}

func (settings *Settings) Validator() *ValidatorSettings {
	return settings.validator
}

func (settings *Settings) Log() *LogSettings {
	return settings.log
}

func (settings *Settings) ProtocolBytes() []byte {
	return settings.protocol.Bytes()
}
