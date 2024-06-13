package customerrors

type IncompleteWebhookDataError struct {
	Message string
}

func (i *IncompleteWebhookDataError) Error() string {
	return i.Message
}
