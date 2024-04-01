package p2p

import (
	"github.com/my-cloud/ruthenium/validatornode/presentation/network"
)

type NeighborsManager interface {
	AddTargets(targets []string)
	HostTarget() string
	Incentive(target string)
	Neighbors() []network.NeighborController
}
