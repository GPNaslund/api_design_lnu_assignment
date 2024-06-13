package customerrors

type IncompleteDogDataError struct {
	Message string
}

func (i *IncompleteDogDataError) Error() string {
	return i.Message
}
