package network

type Synchronizable interface {
	Synchronize()
	Neighbors() []Requestable
}
