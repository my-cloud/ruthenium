package config

type Settings struct {
	HalfLifeInDays              float64
	IncomeBaseInParticles       uint64
	IncomeLimitInParticles      uint64
	MinimalTransactionFee       uint64
	ParticlesPerToken           uint64
	ValidationIntervalInSeconds int64
}
