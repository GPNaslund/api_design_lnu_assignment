package webhook

import (
	dogdto "1dv027/aad/internal/dto/dog"
	webhookdto "1dv027/aad/internal/dto/user/webhook"
	"1dv027/aad/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type UserWebhooksRepository interface {
	GetAllWebhooksByAction(ctx context.Context, action model.WebhookAction) ([]model.Webhook, error)
}

type CryptographyService interface {
	DecryptCipherText(cipherText string) (string, error)
}

type WebhookDispatcher struct {
	userWebhookRepo UserWebhooksRepository
	cryptoService   CryptographyService
}

func NewWebhookDispatcher(userWebhookRepo UserWebhooksRepository, cryptoService CryptographyService) WebhookDispatcher {
	return WebhookDispatcher{
		userWebhookRepo: userWebhookRepo,
		cryptoService:   cryptoService,
	}
}

func (w WebhookDispatcher) DispatchNewDogWebhook(ctx context.Context, dogData dogdto.DogDTO) {
	go func() {
		allWebhooks, err := w.userWebhookRepo.GetAllWebhooksByAction(ctx, model.NEW_DOG_ADDED)
		if err != nil {
			log.Printf("Failed to retrieve webhooks: %v", err)
			return
		}
		for _, userWebhook := range allWebhooks {
			decryptedSecret, err := w.cryptoService.DecryptCipherText(userWebhook.ClientSecret)
			if err != nil {
				log.Printf("Failed to decrypt client secret: %v", err)
				continue
			}
			retries := 3
			for attempt := 0; attempt < retries; attempt++ {
				payLoad := webhookdto.NewDogDispatchDTO{
					NewDog: dogData,
					Secret: decryptedSecret,
				}
				data, err := json.Marshal(payLoad)
				if err != nil {
					log.Printf("Error marshaling payload: %v", err)
					break
				}

				req, err := http.NewRequest("POST", userWebhook.EndpointUrl, bytes.NewBuffer(data))
				if err != nil {
					log.Printf("Error creating HTTP request: %v", err)
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				resp, err := http.DefaultClient.Do(req)
				if err != nil || (resp != nil && resp.StatusCode >= 400) {
					log.Printf("Error dispatching webhook to %s, attempt %d: %v", userWebhook.EndpointUrl, attempt+1, err)
					if resp != nil {
						resp.Body.Close()
					}
					time.Sleep(2 * time.Second)
					continue
				}

				if resp != nil {
					resp.Body.Close()
				}
				log.Printf("Successfully dispatched webhook to %s", userWebhook.EndpointUrl)
				break
			}
		}
	}()
}
