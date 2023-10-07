package protocol

import "time"

type Settings interface {
	BlocksCountLimit() uint64
	GenesisAmountInParticles() uint64
	HalfLifeInNanoseconds() float64
	IncomeBaseInParticles() uint64
	IncomeLimitInParticles() uint64
	MinimalTransactionFee() uint64
	ValidationTimeout() time.Duration
	ValidationTimestamp() int64
}
