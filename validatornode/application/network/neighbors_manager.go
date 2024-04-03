package network

import "github.com/my-cloud/ruthenium/validatornode/presentation"

type NeighborsManager interface {
	AddTargets(targets []string)
	HostTarget() string
	Incentive(target string)
	Neighbors() []presentation.NeighborCaller
}
