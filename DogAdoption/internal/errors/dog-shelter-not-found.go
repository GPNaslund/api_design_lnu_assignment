package customerrors

type DogShelterNotFoundError struct {
	Message string
}

func (d *DogShelterNotFoundError) Error() string {
	return d.Message
}
