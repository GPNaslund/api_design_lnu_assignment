package customerrors

type AdminNotFoundError struct {
	Message string
}

func (a *AdminNotFoundError) Error() string {
	return a.Message
}
