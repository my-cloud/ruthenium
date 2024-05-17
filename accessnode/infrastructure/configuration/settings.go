package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/configuration"
	"io"
	"os"
)

type settingsDto struct {
	Host      *HostSettings
	Template  *TemplateSettings
	Validator *configuration.HostSettings
	Log       *configuration.LogSettings
}

type Settings struct {
	host      *HostSettings
	template  *TemplateSettings
	validator *configuration.HostSettings
	log       *configuration.LogSettings
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
	settings.validator = dto.Validator
	settings.template = dto.Template
	settings.log = dto.Log
	return nil
}

func (settings *Settings) Host() *HostSettings {
	return settings.host
}

func (settings *Settings) Template() *TemplateSettings {
	return settings.template
}

func (settings *Settings) Validator() *configuration.HostSettings {
	return settings.validator
}

func (settings *Settings) Log() *configuration.LogSettings {
	return settings.log
}
