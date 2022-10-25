package network

import "github.com/my-cloud/ruthenium/src/api/node"

type Synchronizable interface {
	Synchronize()
	Neighbors() []Requestable
	AddTargets(request []node.TargetRequest)
	StartSynchronization()
	Wait()
}
