package protocol

type Validator interface {
	Start()
	Stop()
	Do()
}
