package customerrors

type InvalidWebhookDataError struct {
	Message string
}

func (i *InvalidWebhookDataError) Error() string {
	return i.Message
}
