package customerrors

type IntegerConversionError struct {
	Message string
}

func (i *IntegerConversionError) Error() string {
	return i.Message
}
