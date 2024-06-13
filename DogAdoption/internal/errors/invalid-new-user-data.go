package customerrors

type InvalidNewUserDataError struct {
	Message string
}

func (i *InvalidNewUserDataError) Error() string {
	return i.Message
}
