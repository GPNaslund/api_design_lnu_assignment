package model

type Webhook struct {
	Id           int             `json:"id"`
	EndpointUrl  string          `json:"endpoint_url"`
	ClientSecret string          `json:"client_secret"`
	Actions      []WebhookAction `json:"webhook_actions"`
	UserId       int             `json:"user_id"`
}

func (w Webhook) ToJson() map[string]any {
	return map[string]any{
		"id":              w.Id,
		"endpoint_url":    w.EndpointUrl,
		"client_secret":   w.ClientSecret,
		"webhook_actions": w.Actions,
		"user_id":         w.UserId,
	}
}

type WebhookAction string

const (
	NEW_DOG_ADDED WebhookAction = "new_dog_added"
)
