package customerrors

type DogFromJSONError struct {
	Message string
}

func (d *DogFromJSONError) Error() string {
	return d.Message
}
