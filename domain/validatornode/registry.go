package validatornode

type Registry interface {
	IsRegistered(address string) (bool, error)
}
