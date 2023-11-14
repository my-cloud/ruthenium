package network

type Synchronizer interface {
	AddTargets(targets []string)
	HostTarget() string
	Incentive(target string)
	Neighbors() []Neighbor
}
