package protocol

type InputInfoProvider interface {
	TransactionId() string
	OutputIndex() uint16
}
