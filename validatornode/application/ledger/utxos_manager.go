package ledger

import "github.com/my-cloud/ruthenium/validatornode/domain/protocol"

type UtxosManager interface {
	CalculateFee(transaction *protocol.Transaction, timestamp int64) (uint64, error)
	Clear()
	Copy() UtxosManager
	UpdateUtxos(transactions []*protocol.Transaction, timestamp int64) error
	Utxos(address string) []*protocol.Utxo
}
