package configuration

import (
	"encoding/json"
)

type templateSettingsDto struct {
	Path string
}

type TemplateSettings struct {
	path string
}

func (settings *TemplateSettings) UnmarshalJSON(data []byte) error {
	var dto *templateSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.path = dto.Path
	return nil
}

func (settings *TemplateSettings) Path() string {
	return settings.path
}
