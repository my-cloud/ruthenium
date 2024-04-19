package application

import "github.com/my-cloud/ruthenium/validatornode/domain/ledger"

type UtxosManager interface {
	CalculateFee(transaction *ledger.Transaction, timestamp int64) (uint64, error)
	Clear()
	Copy() UtxosManager
	UpdateUtxos(transactions []*ledger.Transaction, timestamp int64) error
	Utxos(address string) []*ledger.Utxo
}
