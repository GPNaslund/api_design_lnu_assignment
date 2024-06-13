package customerrors

type UnauthorizedError struct {
	Message string
}

func (u *UnauthorizedError) Error() string {
	return u.Message
}
