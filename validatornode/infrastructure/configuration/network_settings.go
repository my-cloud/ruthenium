package configuration

import (
	"encoding/json"
	"time"
)

type networkSettingsDto struct {
	ConnectionTimeoutInSeconds       int
	MaxOutboundsCount                int
	Seeds                            []string
	SynchronizationIntervalInSeconds int
}

type NetworkSettings struct {
	connectionTimeout    time.Duration
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
	settings.connectionTimeout = time.Duration(dto.ConnectionTimeoutInSeconds) * time.Second
	settings.maxOutboundsCount = dto.MaxOutboundsCount
	settings.seeds = dto.Seeds
	settings.synchronizationTimer = time.Duration(dto.SynchronizationIntervalInSeconds) * time.Second
	return nil
}

func (settings *NetworkSettings) ConnectionTimeout() time.Duration {
	return settings.connectionTimeout
}

func (settings *NetworkSettings) MaxOutboundsCount() int {
	return settings.maxOutboundsCount
}

func (settings *NetworkSettings) SynchronizationTimer() time.Duration {
	return settings.synchronizationTimer
}

func (settings *NetworkSettings) Seeds() map[string]int {
	scoresBySeedTargetValue := map[string]int{}
	for _, seedStringTargetValue := range settings.seeds {
		scoresBySeedTargetValue[seedStringTargetValue] = 0
	}
	return scoresBySeedTargetValue
}
