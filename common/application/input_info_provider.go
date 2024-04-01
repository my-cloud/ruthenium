package application

type InputInfoProvider interface {
	TransactionId() string
	OutputIndex() uint16
}
