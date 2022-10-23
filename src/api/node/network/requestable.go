package network

type Requestable interface {
	Ip() string
	Port() uint16
	Target() string
	IsFound() bool
	GetBlocks() (blockResponses []*BlockResponse, err error)
	SendTargets(request []TargetRequest) (err error)
	AddTransaction(request TransactionRequest) (err error)
	GetTransactions() (transactionResponses []TransactionResponse, err error)
	GetAmount(request AmountRequest) (amountResponse *AmountResponse, err error)
	Mine() (err error)
	StartMining() (err error)
	StopMining() error
}
