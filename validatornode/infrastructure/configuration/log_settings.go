package configuration

import (
	"encoding/json"
)

type logSettingsDto struct {
	Level string
}

type LogSettings struct {
	level string
}

func (settings *LogSettings) UnmarshalJSON(data []byte) error {
	var dto *logSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.level = dto.Level
	return nil
}

func (settings *LogSettings) Level() string {
	return settings.level
}
