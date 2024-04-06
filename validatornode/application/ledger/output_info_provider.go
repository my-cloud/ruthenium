package ledger

type OutputInfoProvider interface {
	InitialValue() uint64
	IsYielding() bool
}
