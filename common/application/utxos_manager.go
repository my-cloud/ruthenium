package application

type UtxosManager interface {
	Clear()
	Copy() UtxosManager
	UpdateUtxos(transactionsBytes []byte, timestamp int64) error
	Utxo(input InputInfoProvider) (UtxoInfoProvider, error)
	Utxos(address string) []byte
}
