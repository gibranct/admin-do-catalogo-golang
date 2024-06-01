package domain

type NotFoundException struct {
	Message string
}

func (ex NotFoundException) Error() string {
	return ex.Message
}
