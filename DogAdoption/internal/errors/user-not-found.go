package customerrors

type UserNotFoundError struct {
	Message string
}

func (u *UserNotFoundError) Error() string {
	return u.Message
}
