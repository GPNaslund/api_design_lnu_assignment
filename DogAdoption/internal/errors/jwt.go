package customerrors

type JwtError struct {
	Message string
}

func (j *JwtError) Error() string {
	return j.Message
}
