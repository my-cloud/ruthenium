package network

type Validator interface {
	StartValidation()
	StopValidation()
	Validate()
}
