package presentation

type Pulser interface {
	Start()
	Stop()
	Pulse()
}
