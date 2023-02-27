package pohtest

type Registry struct{}

func NewRegistry() *Registry { return &Registry{} }

func (registry *Registry) IsRegistered(string) (isRegistered bool, err error) { return true, nil }
