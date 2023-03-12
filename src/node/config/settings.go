package config

type Settings struct {
	GenesisAmountInParticles         uint64
	MaxOutboundsCount                int
	MinimalTransactionFee            uint64
	SynchronizationIntervalInSeconds int
	ValidationIntervalInSeconds      int
	VerificationsCountPerValidation  int
}
