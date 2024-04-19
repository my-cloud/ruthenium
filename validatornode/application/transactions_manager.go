package application

import "github.com/my-cloud/ruthenium/validatornode/domain/ledger"

type TransactionsManager interface {
	AddTransaction(transaction *ledger.Transaction, broadcasterTarget string, hostTarget string)
	Transactions() []*ledger.Transaction
}
