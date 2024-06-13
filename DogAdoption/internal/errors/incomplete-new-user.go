package customerrors

type IncompleteNewUserError struct {
	Message string
}

func (i *IncompleteNewUserError) Error() string {
	return i.Message
}
