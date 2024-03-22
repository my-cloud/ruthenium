package ledger

import (
	"github.com/my-cloud/ruthenium/domain"
)

type UtxosManager interface {
	Clear()
	Copy() UtxosManager
	UpdateUtxos(transactions []*Transaction, timestamp int64) error // TODO remove dependency to Transaction and move to validation package
	Utxo(input domain.InputInfoProvider) (domain.UtxoInfoProvider, error)
	Utxos(address string) []byte
	VerifyUtxos(transactions []*Transaction, timestamp int64) error
}
