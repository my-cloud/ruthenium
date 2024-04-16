package application

import "github.com/my-cloud/ruthenium/validatornode/domain/protocol"

type TransactionsManager interface {
	AddTransaction(transaction *protocol.Transaction, broadcasterTarget string, hostTarget string)
	Transactions() []*protocol.Transaction
}
