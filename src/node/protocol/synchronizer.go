package protocol

import "github.com/my-cloud/ruthenium/src/node/neighborhood"

type Synchronizer interface {
	Neighbors() []neighborhood.Neighbor
}
