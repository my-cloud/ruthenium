package network

type Synchronizable interface {
	Synchronize()
	Neighbors() []Requestable
	AddTargets(request []TargetRequest)
	StartSynchronization()
	Wait()
}
