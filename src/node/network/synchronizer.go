package network

type Synchronizer interface {
	AddTargets(requests []TargetRequest)
	HostTarget() string
	Incentive(target string)
	Neighbors() []Neighbor
}
