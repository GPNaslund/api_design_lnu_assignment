package customerrors

type DatabaseError struct {
	Message string
}

func (d *DatabaseError) Error() string {
	return d.Message
}
