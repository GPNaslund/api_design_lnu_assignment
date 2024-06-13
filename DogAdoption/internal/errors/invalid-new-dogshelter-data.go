package customerrors

type InvalidNewDogShelterDataError struct {
	Message string
}

func (i *InvalidNewDogShelterDataError) Error() string {
	return i.Message
}
