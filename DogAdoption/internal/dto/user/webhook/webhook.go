package userwebhookdto

import "1dv027/aad/internal/model"

type UserWebhookDTO struct {
	EndpointUrl string                `json:"endpoint_url"`
	Actions     []model.WebhookAction `json:"webhook_actions"`
	UserId      int                   `json:"user_id"`
}
