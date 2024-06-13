package customerrors

type IncompleteDogShelterDataError struct {
	Message string
}

func (i *IncompleteDogShelterDataError) Error() string {
	return i.Message
}
