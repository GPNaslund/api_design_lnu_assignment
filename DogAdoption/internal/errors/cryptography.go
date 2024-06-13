package customerrors

type CryptographyError struct {
	Message string
}

func (c *CryptographyError) Error() string {
	return c.Message
}
