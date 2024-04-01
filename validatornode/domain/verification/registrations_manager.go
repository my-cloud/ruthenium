package verification

type RegistrationsManager interface {
	Clear()
	Copy() RegistrationsManager
	Filter(addresses []string) (newAddresses []string)
	IsRegistered(address string) bool
	RemovedAddresses() (removedAddresses []string)
	Update(addedAddresses []string, removedAddresses []string)
	Verify(addedAddresses []string, removedAddresses []string) error
}
