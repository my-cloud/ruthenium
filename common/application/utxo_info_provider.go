package application

type UtxoInfoProvider interface {
	Address() string
	Value(currentTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64) uint64
}
