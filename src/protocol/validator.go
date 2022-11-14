package protocol

type Validator interface {
	StartValidation()
	StopValidation()
	Validate()
}
