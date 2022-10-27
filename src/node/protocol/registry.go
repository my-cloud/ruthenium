package protocol

type Registry interface {
	IsRegistered(address string) (bool, error)
}
