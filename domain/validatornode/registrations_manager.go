package validatornode

type RegistrationsManager interface {
	IsRegistered(address string) (bool, error)
}
