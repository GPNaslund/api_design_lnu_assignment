package customerrors

type DogNotFoundError struct {
	Message string
}

func (d *DogNotFoundError) Error() string {
	return d.Message
}
