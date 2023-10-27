package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type settingsDto struct {
	BlocksCountLimit                 uint64
	GenesisAmountInParticles         uint64
	HalfLifeInDays                   float64
	IncomeBaseInParticles            uint64
	IncomeLimitInParticles           uint64
	MaxOutboundsCount                int
	MinimalTransactionFee            uint64
	ParticlesPerToken                uint64
	SynchronizationIntervalInSeconds int
	ValidationIntervalInSeconds      int64
	ValidationTimeoutInSeconds       int64
	VerificationsCountPerValidation  int64
}

type Settings struct {
	bytes                           []byte
	blocksCountLimit                uint64
	genesisAmountInParticles        uint64
	halfLifeInNanoseconds           float64
	incomeBaseInParticles           uint64
	incomeLimitInParticles          uint64
	maxOutboundsCount               int
	minimalTransactionFee           uint64
	particlesPerToken               uint64
	synchronizationTimer            time.Duration
	validationTimestamp             int64
	validationTimer                 time.Duration
	validationTimeout               time.Duration
	verificationsCountPerValidation int64
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
	settings.bytes = data
	settings.blocksCountLimit = dto.BlocksCountLimit
	settings.genesisAmountInParticles = dto.GenesisAmountInParticles
	hoursByDay := 24.
	settings.halfLifeInNanoseconds = dto.HalfLifeInDays * hoursByDay * float64(time.Hour.Nanoseconds())
	settings.incomeBaseInParticles = dto.IncomeBaseInParticles
	settings.incomeLimitInParticles = dto.IncomeLimitInParticles
	settings.maxOutboundsCount = dto.MaxOutboundsCount
	settings.minimalTransactionFee = dto.MinimalTransactionFee
	settings.particlesPerToken = dto.ParticlesPerToken
	settings.synchronizationTimer = time.Duration(dto.SynchronizationIntervalInSeconds) * time.Second
	settings.validationTimestamp = dto.ValidationIntervalInSeconds * time.Second.Nanoseconds()
	settings.validationTimer = time.Duration(dto.ValidationIntervalInSeconds) * time.Second
	settings.validationTimeout = time.Duration(dto.ValidationTimeoutInSeconds) * time.Second
	settings.verificationsCountPerValidation = dto.VerificationsCountPerValidation
	return nil
}

func (settings *Settings) Bytes() []byte {
	return settings.bytes
}

func (settings *Settings) BlocksCountLimit() uint64 {
	return settings.blocksCountLimit
}

func (settings *Settings) GenesisAmountInParticles() uint64 {
	return settings.genesisAmountInParticles
}

func (settings *Settings) HalfLifeInNanoseconds() float64 {
	return settings.halfLifeInNanoseconds
}

func (settings *Settings) IncomeBaseInParticles() uint64 {
	return settings.incomeBaseInParticles
}

func (settings *Settings) IncomeLimitInParticles() uint64 {
	return settings.incomeLimitInParticles
}

func (settings *Settings) MaxOutboundsCount() int {
	return settings.maxOutboundsCount
}

func (settings *Settings) MinimalTransactionFee() uint64 {
	return settings.minimalTransactionFee
}

func (settings *Settings) ParticlesPerToken() uint64 {
	return settings.particlesPerToken
}

func (settings *Settings) SynchronizationTimer() time.Duration {
	return settings.synchronizationTimer
}

func (settings *Settings) ValidationTimer() time.Duration {
	return settings.validationTimer
}

func (settings *Settings) ValidationTimestamp() int64 {
	return settings.validationTimestamp
}

func (settings *Settings) ValidationTimeout() time.Duration {
	return settings.validationTimeout
}

func (settings *Settings) VerificationsCountPerValidation() int64 {
	return settings.verificationsCountPerValidation
}
