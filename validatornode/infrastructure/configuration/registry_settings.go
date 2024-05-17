package configuration

import (
	"encoding/json"
	"time"
)

type registrySettingsDto struct {
	SynchronizationIntervalInSeconds int
}

type RegistrySettings struct {
	synchronizationTimer time.Duration
}

func (settings *RegistrySettings) UnmarshalJSON(data []byte) error {
	var dto *registrySettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.synchronizationTimer = time.Duration(dto.SynchronizationIntervalInSeconds) * time.Second
	return nil
}

func (settings *RegistrySettings) SynchronizationTimer() time.Duration {
	return settings.synchronizationTimer
}
