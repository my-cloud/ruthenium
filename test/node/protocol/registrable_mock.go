package protocol

type RegistrableMock struct {
}

func NewRegistrableMock() *RegistrableMock {
	return &RegistrableMock{}
}

var IsRegisteredMock func(string) (bool, error)

func (mock *RegistrableMock) IsRegistered(address string) (bool, error) {
	return IsRegisteredMock(address)
}
