package userwebhookdto

type UpdateUserWebhookDTO struct {
	EndpointUrl  *string   `json:"endpoint_url"`
	Actions      *[]string `json:"webhook_actions"`
	ClientSecret *string   `json:"client_secret"`
}
