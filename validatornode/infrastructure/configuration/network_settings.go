package configuration

import (
	"encoding/json"
	"time"
)

type networkSettingsDto struct {
	MaxOutboundsCount                int
	Seeds                            []string
	SynchronizationIntervalInSeconds int
}

type NetworkSettings struct {
	maxOutboundsCount    int
	seeds                []string
	synchronizationTimer time.Duration
}

func (settings *NetworkSettings) UnmarshalJSON(data []byte) error {
	var dto *networkSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.maxOutboundsCount = dto.MaxOutboundsCount
	settings.seeds = dto.Seeds
	settings.synchronizationTimer = time.Duration(dto.SynchronizationIntervalInSeconds) * time.Second
	return nil
}

func (settings *NetworkSettings) MaxOutboundsCount() int {
	return settings.maxOutboundsCount
}

func (settings *NetworkSettings) SynchronizationTimer() time.Duration {
	return settings.synchronizationTimer
}

func (settings *NetworkSettings) Seeds() []string {
	return settings.seeds
}
