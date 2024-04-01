package application

type Pulser interface {
	Start()
	Stop()
	Pulse()
}
