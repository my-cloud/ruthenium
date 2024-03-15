package domain

type Pulser interface {
	Start()
	Stop()
	Pulse()
}
