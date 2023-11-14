package observernode

type Settings interface {
	HalfLifeInNanoseconds() float64
	IncomeBase() uint64
	IncomeLimit() uint64
	MinimalTransactionFee() uint64
	SmallestUnitsPerCoin() uint64
	ValidationTimestamp() int64
}
