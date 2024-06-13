package userwebhookdto

type NewUserWebhookDTO struct {
	EndpointUrl  *string   `json:"endpoint_url"`
	Actions      *[]string `json:"webhook_actions"`
	ClientSecret *string   `json:"client_secret"`
}
