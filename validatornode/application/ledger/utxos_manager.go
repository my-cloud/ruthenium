package ledger

type UtxosManager interface {
	CalculateFee(inputs []InputInfoProvider, outputs []OutputInfoProvider, timestamp int64) (uint64, error)
	Clear()
	Copy() UtxosManager
	UpdateUtxos(transactionsBytes []byte, timestamp int64) error
	Utxos(address string) []byte
}
