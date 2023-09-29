package protocol

type Utxo interface {
	Value(currentTimestamp int64, genesisTimestamp int64, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64, validationTimestamp int64) uint64
}
