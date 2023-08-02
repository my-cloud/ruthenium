package config

type Settings struct {
	GenesisAmountInParticles         uint64
	HalfLifeInDays                   float64
	MaxOutboundsCount                int
	MinimalTransactionFee            uint64
	SynchronizationIntervalInSeconds int
	ValidationIntervalInSeconds      int64
	VerificationsCountPerValidation  int
}
