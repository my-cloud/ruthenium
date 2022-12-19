package network

type Synchronizer interface {
	AddTargets(requests []TargetRequest)
	Neighbors() []Neighbor
}
