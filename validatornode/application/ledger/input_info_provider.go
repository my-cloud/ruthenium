package ledger

type InputInfoProvider interface {
	TransactionId() string
	OutputIndex() uint16
}
