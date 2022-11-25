package protocol

type Engine interface {
	Start()
	Stop()
	Do()
}
