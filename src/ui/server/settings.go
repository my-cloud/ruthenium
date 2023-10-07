package server

type Settings interface {
	HalfLifeInNanoseconds() float64
	IncomeBaseInParticles() uint64
	IncomeLimitInParticles() uint64
	MinimalTransactionFee() uint64
	ParticlesPerToken() uint64
	ValidationTimestamp() int64
}
