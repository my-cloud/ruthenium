package domain

type InputInfo interface {
	TransactionId() string
	OutputIndex() uint16
}
