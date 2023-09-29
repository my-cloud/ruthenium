package protocol

type Input interface {
	TransactionId() string
	OutputIndex() uint16
}
