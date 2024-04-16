package wallet

type SettingsProvider interface {
	HalfLifeInNanoseconds() float64
	IncomeBase() uint64
	IncomeLimit() uint64
	SmallestUnitsPerCoin() uint64
}
