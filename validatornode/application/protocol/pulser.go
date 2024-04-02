package protocol

type Pulser interface {
	Start()
	Stop()
	Pulse()
}
