package ledger

type TransactionsManager interface {
	AddTransaction(transactionRequestBytes []byte, hostTarget string)
	Transactions() []byte
}
