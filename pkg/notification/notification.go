package notification

import "github.com.br/gibranct/admin-do-catalogo/pkg/validator"

type Notification struct {
	errors []error
}

func (n *Notification) Add(err error) validator.ValidationHandler {
	n.errors = append(n.errors, err)
	return n
}

func (n Notification) GetErrors() []error {
	return n.errors
}

func (n Notification) HasErrors() bool {
	return len(n.errors) != 0
}

func CreateNotification() *Notification {
	return &Notification{
		errors: make([]error, 0, 10),
	}
}
