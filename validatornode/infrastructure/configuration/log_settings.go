package configuration

import (
	"encoding/json"
)

type logSettingsDto struct {
	LogLevel string
}

type LogSettings struct {
	logLevel string
}

func (settings *LogSettings) UnmarshalJSON(data []byte) error {
	var dto *logSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.logLevel = dto.LogLevel
	return nil
}

func (settings *LogSettings) LogLevel() string {
	return settings.logLevel
}
