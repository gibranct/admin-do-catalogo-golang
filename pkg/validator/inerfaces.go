package validator

type ValidationHandler interface {
	HasErrors() bool
	GetErrors() []error
	Add(err error) ValidationHandler
}

type Validator interface {
	Validate()
	ValidationHandler
}

type Entity interface {
	Validate(handler ValidationHandler)
}
