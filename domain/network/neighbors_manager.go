package network

type NeighborsManager interface {
	AddTargets(targets []string)
	HostTarget() string
	Incentive(target string)
	Neighbors() []Neighbor
}
