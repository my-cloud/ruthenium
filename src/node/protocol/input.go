package protocol

type InputInfo interface {
	TransactionId() string
	OutputIndex() uint16
}
