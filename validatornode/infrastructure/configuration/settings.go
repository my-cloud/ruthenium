package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type settingsDto struct {
	Host struct {
		Ip   string
		Port int
	}
	Network struct {
		MaxOutboundsCount                int
		Seeds                            []string
		SynchronizationIntervalInSeconds int
		ValidationTimeoutInSeconds       int64
	}
	Protocol struct {
		BlocksCountLimit                uint64
		GenesisAmount                   uint64
		HalfLifeInDays                  float64
		IncomeBase                      uint64
		IncomeLimit                     uint64
		MinimalTransactionFee           uint64
		SmallestUnitsPerCoin            uint64
		ValidationIntervalInSeconds     int64
		VerificationsCountPerValidation int64
	}
	Validator struct {
		Address   string
		InfuraKey string
	}
	Log struct {
		LogLevel string
	}
}

type Settings struct {
	bytes []byte
	host  struct {
		ip   string
		port string
	}
	network struct {
		maxOutboundsCount    int
		seeds                []string
		synchronizationTimer time.Duration
		validationTimeout    time.Duration
	}
	protocol struct {
		blocksCountLimit                uint64
		genesisAmount                   uint64
		halfLifeInNanoseconds           float64
		incomeBase                      uint64
		incomeLimit                     uint64
		minimalTransactionFee           uint64
		validationTimer                 time.Duration
		smallestUnitsPerCoin            uint64
		validationTimestamp             int64
		verificationsCountPerValidation int64
	}
	validator struct {
		address   string
		infuraKey string
	}
	log struct {
		logLevel string
	}
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
	settings.host.ip = dto.Host.Ip
	settings.host.port = strconv.Itoa(dto.Host.Port)
	settings.network.maxOutboundsCount = dto.Network.MaxOutboundsCount
	settings.network.seeds = dto.Network.Seeds
	settings.network.synchronizationTimer = time.Duration(dto.Network.SynchronizationIntervalInSeconds) * time.Second
	settings.network.validationTimeout = time.Duration(dto.Network.ValidationTimeoutInSeconds) * time.Second
	settings.protocol.blocksCountLimit = dto.Protocol.BlocksCountLimit
	settings.protocol.genesisAmount = dto.Protocol.GenesisAmount
	hoursByDay := 24.
	settings.protocol.halfLifeInNanoseconds = dto.Protocol.HalfLifeInDays * hoursByDay * float64(time.Hour.Nanoseconds())
	settings.protocol.incomeBase = dto.Protocol.IncomeBase
	settings.protocol.incomeLimit = dto.Protocol.IncomeLimit
	settings.protocol.minimalTransactionFee = dto.Protocol.MinimalTransactionFee
	settings.protocol.validationTimestamp = dto.Protocol.ValidationIntervalInSeconds * time.Second.Nanoseconds()
	settings.protocol.validationTimer = time.Duration(dto.Protocol.ValidationIntervalInSeconds) * time.Second
	settings.protocol.verificationsCountPerValidation = dto.Protocol.VerificationsCountPerValidation
	settings.protocol.smallestUnitsPerCoin = dto.Protocol.SmallestUnitsPerCoin
	settings.validator.address = dto.Validator.Address
	settings.validator.infuraKey = dto.Validator.InfuraKey
	settings.log.logLevel = dto.Log.LogLevel
	return nil
}

func (settings *Settings) Bytes() []byte {
	return settings.bytes
}

func (settings *Settings) BlocksCountLimit() uint64 {
	return settings.protocol.blocksCountLimit
}

func (settings *Settings) GenesisAmount() uint64 {
	return settings.protocol.genesisAmount
}

func (settings *Settings) HalfLifeInNanoseconds() float64 {
	return settings.protocol.halfLifeInNanoseconds
}

func (settings *Settings) IncomeBase() uint64 {
	return settings.protocol.incomeBase
}

func (settings *Settings) IncomeLimit() uint64 {
	return settings.protocol.incomeLimit
}

func (settings *Settings) MaxOutboundsCount() int {
	return settings.network.maxOutboundsCount
}

func (settings *Settings) MinimalTransactionFee() uint64 {
	return settings.protocol.minimalTransactionFee
}

func (settings *Settings) SmallestUnitsPerCoin() uint64 {
	return settings.protocol.smallestUnitsPerCoin
}

func (settings *Settings) SynchronizationTimer() time.Duration {
	return settings.network.synchronizationTimer
}

func (settings *Settings) ValidationTimer() time.Duration {
	return settings.protocol.validationTimer
}

func (settings *Settings) ValidationTimestamp() int64 {
	return settings.protocol.validationTimestamp
}

func (settings *Settings) ValidationTimeout() time.Duration {
	return settings.network.validationTimeout
}

func (settings *Settings) VerificationsCountPerValidation() int64 {
	return settings.protocol.verificationsCountPerValidation
}

func (settings *Settings) LogLevel() string {
	return settings.log.logLevel
}

func (settings *Settings) InfuraKey() string {
	return settings.validator.infuraKey
}

func (settings *Settings) Address() string {
	return settings.validator.address
}

func (settings *Settings) Seeds() []string {
	return settings.network.seeds
}

func (settings *Settings) Ip() string {
	return settings.host.ip
}

func (settings *Settings) Port() string {
	return settings.host.port
}
