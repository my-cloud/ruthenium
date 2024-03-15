package domain

type InputInfoProvider interface {
	TransactionId() string
	OutputIndex() uint16
}
