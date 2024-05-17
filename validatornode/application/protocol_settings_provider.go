package application

import "time"

type ProtocolSettingsProvider interface {
	BlocksCountLimit() uint64
	GenesisAmount() uint64
	HalfLifeInNanoseconds() float64
	IncomeBase() uint64
	IncomeLimit() uint64
	MinimalTransactionFee() uint64
	SmallestUnitsPerCoin() uint64
	ValidationTimeout() time.Duration
	ValidationTimer() time.Duration
	ValidationTimestamp() int64
	VerificationsCountPerValidation() int64
}
