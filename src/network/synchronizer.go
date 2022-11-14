package network

type Synchronizer interface {
	Neighbors() []Neighbor
}
