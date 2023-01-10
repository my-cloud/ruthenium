package network

type Neighbor interface {
	Ip() string
	Port() uint16
	Target() string
	GetBlocks() (blockResponses []*BlockResponse, err error)
	SendTargets(request []TargetRequest) (err error)
	AddTransaction(request TransactionRequest) (err error)
	GetTransactions() (transactionResponses []TransactionResponse, err error)
	GetAmount(request AmountRequest) (amountResponse *AmountResponse, err error)
	StartMining() (err error)
	StopMining() error
}
