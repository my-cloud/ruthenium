package validatornode

import "time"

type SettingsProvider interface {
	BlocksCountLimit() uint64
	GenesisAmount() uint64
	HalfLifeInNanoseconds() float64
	IncomeBase() uint64
	IncomeLimit() uint64
	MinimalTransactionFee() uint64
	ValidationTimeout() time.Duration
	ValidationTimestamp() int64
}
