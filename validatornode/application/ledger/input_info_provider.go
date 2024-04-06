package ledger

type InputInfoProvider interface {
	TransactionId() string
	OutputIndex() uint16
	Address() string
}
