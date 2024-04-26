package configuration

import (
	"encoding/json"
	"math"
	"time"
)

type protocolSettingsDto struct {
	BlocksCountLimit                uint64
	CoinDigitsCount                 uint8
	GenesisAmount                   uint64
	HalfLifeInDays                  float64
	IncomeBase                      uint64
	IncomeLimit                     uint64
	MinimalTransactionFee           uint64
	ValidationIntervalInSeconds     int64
	ValidationTimeoutInSeconds      int64
	VerificationsCountPerValidation int64
}

type ProtocolSettings struct {
	bytes                           []byte
	blocksCountLimit                uint64
	genesisAmount                   uint64
	halfLifeInNanoseconds           float64
	incomeBase                      uint64
	incomeLimit                     uint64
	minimalTransactionFee           uint64
	smallestUnitsPerCoin            uint64
	validationTimeout               time.Duration
	validationTimer                 time.Duration
	validationTimestamp             int64
	verificationsCountPerValidation int64
}

func (settings *ProtocolSettings) UnmarshalJSON(data []byte) error {
	var dto *protocolSettingsDto
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}
	settings.bytes = data
	settings.blocksCountLimit = dto.BlocksCountLimit
	settings.genesisAmount = dto.GenesisAmount
	hoursByDay := 24.
	settings.halfLifeInNanoseconds = dto.HalfLifeInDays * hoursByDay * float64(time.Hour.Nanoseconds())
	settings.incomeBase = dto.IncomeBase
	settings.incomeLimit = dto.IncomeLimit
	settings.minimalTransactionFee = dto.MinimalTransactionFee
	settings.smallestUnitsPerCoin = uint64(math.Pow10(int(dto.CoinDigitsCount)))
	settings.validationTimeout = time.Duration(dto.ValidationTimeoutInSeconds) * time.Second
	settings.validationTimer = time.Duration(dto.ValidationIntervalInSeconds) * time.Second
	settings.validationTimestamp = dto.ValidationIntervalInSeconds * time.Second.Nanoseconds()
	settings.verificationsCountPerValidation = dto.VerificationsCountPerValidation
	return nil
}

func (settings *ProtocolSettings) Bytes() []byte {
	return settings.bytes
}

func (settings *ProtocolSettings) BlocksCountLimit() uint64 {
	return settings.blocksCountLimit
}

func (settings *ProtocolSettings) GenesisAmount() uint64 {
	return settings.genesisAmount
}

func (settings *ProtocolSettings) HalfLifeInNanoseconds() float64 {
	return settings.halfLifeInNanoseconds
}

func (settings *ProtocolSettings) IncomeBase() uint64 {
	return settings.incomeBase
}

func (settings *ProtocolSettings) IncomeLimit() uint64 {
	return settings.incomeLimit
}

func (settings *ProtocolSettings) MinimalTransactionFee() uint64 {
	return settings.minimalTransactionFee
}

func (settings *ProtocolSettings) SmallestUnitsPerCoin() uint64 {
	return settings.smallestUnitsPerCoin
}

func (settings *ProtocolSettings) ValidationTimeout() time.Duration {
	return settings.validationTimeout
}

func (settings *ProtocolSettings) ValidationTimer() time.Duration {
	return settings.validationTimer
}

func (settings *ProtocolSettings) ValidationTimestamp() int64 {
	return settings.validationTimestamp
}

func (settings *ProtocolSettings) VerificationsCountPerValidation() int64 {
	return settings.verificationsCountPerValidation
}
