package network

type Synchronizer interface {
	AddTargets(requests []TargetRequest)
	Incentive(target string)
	Neighbors() []Neighbor
}
