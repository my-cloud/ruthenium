package protocol

type Controllable interface {
	Start()
	Stop()
	Do()
}
