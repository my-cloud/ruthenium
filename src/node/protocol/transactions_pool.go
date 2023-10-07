package protocol

type TransactionsPool interface {
	AddTransaction(transactionRequestBytes []byte, hostTarget string)
	Transactions() []byte
}
