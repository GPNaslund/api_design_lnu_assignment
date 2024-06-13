package userdto

import userwebhookdto "1dv027/aad/internal/dto/user/webhook"

type UserDTO struct {
	Id       int                           `json:"id"`
	Username string                        `json:"username"`
	Webhook  userwebhookdto.UserWebhookDTO `json:"webhook"`
}
