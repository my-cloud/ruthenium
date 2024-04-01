package application

type TransactionsManager interface {
	AddTransaction(transactionRequestBytes []byte, hostTarget string)
	Transactions() []byte
}
