package humanity

type Registrable interface {
	IsRegistered(address string) (bool, error)
}
