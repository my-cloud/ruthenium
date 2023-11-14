package validatornode

type InputInfo interface {
	TransactionId() string
	OutputIndex() uint16
}
