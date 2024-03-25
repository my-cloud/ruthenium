package domain

type UtxosManager interface {
	Clear()
	Copy() UtxosManager
	UpdateUtxos(transactionsBytes []byte, timestamp int64) error // TODO remove dependency to Transaction and move to validation package
	Utxo(input InputInfoProvider) (UtxoInfoProvider, error)
	Utxos(address string) []byte
}
