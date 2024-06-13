package customerrors

type WrongCredentialsError struct {
	Message string
}

func (w *WrongCredentialsError) Error() string {
	return w.Message
}
