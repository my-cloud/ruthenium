package api

type Registrable interface {
	IsRegistered(address string) (bool, error)
}
