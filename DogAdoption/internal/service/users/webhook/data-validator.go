package userwebhookservice

import (
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"fmt"
	"net/url"
	"strings"
)

type WebhookDataValidatorService struct{}

func NewWebhookDataValidatorService() WebhookDataValidatorService {
	return WebhookDataValidatorService{}
}

func (w WebhookDataValidatorService) validateWebhookActions(actions []string) error {
	var errMsgs []string
	for _, webhookAction := range actions {
		switch webhookAction {
		case string(model.NEW_DOG_ADDED):
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("invalid webhook action: %s", webhookAction))
		}
	}

	if len(errMsgs) > 0 {
		errMsgString := strings.Join(errMsgs, " + ")
		return &customerrors.InvalidWebhookDataError{Message: errMsgString}
	}
	return nil
}

func (w WebhookDataValidatorService) validateWebhookEndpointUrl(urlString string) error {
	parsedUrl, err := url.ParseRequestURI(urlString)
	if err != nil {
		return &customerrors.InvalidWebhookDataError{Message: "invalid endpoint url"}
	}
	if parsedUrl.Scheme != "https" {
		return &customerrors.InvalidWebhookDataError{Message: "endpoint url must be https"}
	}
	return nil
}

func (w WebhookDataValidatorService) validateWebhookClientSecret(secret string, minSecretLength int) error {
	if len(secret) < minSecretLength {
		return &customerrors.InvalidWebhookDataError{Message: fmt.Sprintf("secret must be minimum %d charcters long", minSecretLength)}
	}
	return nil
}
