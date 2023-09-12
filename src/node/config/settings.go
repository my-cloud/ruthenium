package config

type Settings struct {
	GenesisAmountInParticles         uint64
	HalfLifeInDays                   float64
	IncomeBaseInParticles            uint64
	IncomeLimitInParticles           uint64
	MaxOutboundsCount                int
	MinimalTransactionFee            uint64
	SynchronizationIntervalInSeconds int
	ValidationIntervalInSeconds      int64
	VerificationsCountPerValidation  int
}
