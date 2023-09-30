package protocol

type TransactionsPool interface {
	AddTransaction(transaction []byte, hostTarget string)
	Transactions() []byte
}
