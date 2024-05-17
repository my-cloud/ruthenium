package verification

type HumansManager interface {
	IsRegistered(address string) (isRegistered bool, err error)
}
