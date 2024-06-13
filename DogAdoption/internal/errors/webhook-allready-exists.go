package customerrors

type WebhookAllreadyExistsError struct {
	Message string
}

func (w *WebhookAllreadyExistsError) Error() string {
	return w.Message
}
