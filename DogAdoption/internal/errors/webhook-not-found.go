package customerrors

type WebhookNotFoundError struct {
	Message string
}

func (w *WebhookNotFoundError) Error() string {
	return w.Message
}
