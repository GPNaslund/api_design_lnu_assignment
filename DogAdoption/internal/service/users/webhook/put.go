package userwebhookservice

import (
	"1dv027/aad/internal/dto"
	webhookdto "1dv027/aad/internal/dto/user/webhook"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"strconv"
)

type PutUserWebhooksRepository interface {
	UpdateUserWebhook(ctx context.Context, userId int, data webhookdto.UpdateUserWebhookDTO) (model.Webhook, error)
}

type PutUserWebhooksDataValidator interface {
	validateWebhookActions(actions []string) error
	validateWebhookEndpointUrl(urlString string) error
	validateWebhookClientSecret(secret string, minSecretLength int) error
}

type PutUserWebhooksCryptographyService interface {
	EncryptPlainText(plainText string) (string, error)
}

type PutUserWebhooksService struct {
	repo          PutUserWebhooksRepository
	dataValidator PutUserWebhooksDataValidator
	cryptoService PutUserWebhooksCryptographyService
}

func NewPutUserWebhookService(repo PutUserWebhooksRepository,
	dataValidator PutUserWebhooksDataValidator, cryptoService PutUserWebhooksCryptographyService) PutUserWebhooksService {
	return PutUserWebhooksService{
		repo:          repo,
		dataValidator: dataValidator,
		cryptoService: cryptoService,
	}
}

func (p PutUserWebhooksService) UpdateUserWebhook(ctx context.Context, idParam string, user dto.UserCredentials,
	data webhookdto.UpdateUserWebhookDTO) (webhookdto.UserWebhookDTO, error) {
	emptyDto := webhookdto.UserWebhookDTO{}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		return emptyDto, &customerrors.IntegerConversionError{}
	}

	if user.UserRole != model.ADMIN && user.UserRole != model.USER {
		return emptyDto, &customerrors.UnauthorizedError{}
	}

	if user.UserRole == model.USER && idParamInt != user.Id {
		return emptyDto, &customerrors.UnauthorizedError{}
	}

	err = p.validateUpdateDtoData(data)
	if err != nil {
		return emptyDto, err
	}

	if data.ClientSecret != nil {
		encryptedClientSecret, err := p.cryptoService.EncryptPlainText(*data.ClientSecret)
		if err != nil {
			return emptyDto, &customerrors.CryptographyError{}
		}
		data.ClientSecret = &encryptedClientSecret
	}

	webhookModel, err := p.repo.UpdateUserWebhook(ctx, idParamInt, data)
	if err != nil {
		return emptyDto, err
	}
	webhookJson, err := json.Marshal(webhookModel.ToJson())
	if err != nil {
		return emptyDto, err
	}
	var webhookDto webhookdto.UserWebhookDTO
	err = json.Unmarshal(webhookJson, &webhookDto)
	if err != nil {
		return emptyDto, err
	}
	return webhookDto, nil
}

func (p PutUserWebhooksService) validateUpdateDtoData(data webhookdto.UpdateUserWebhookDTO) error {
	if data.Actions == nil && data.EndpointUrl == nil && data.ClientSecret == nil {
		return &customerrors.IncompleteWebhookDataError{}
	}

	if data.Actions != nil {
		err := p.dataValidator.validateWebhookActions(*data.Actions)
		if err != nil {
			return err
		}
	}

	if data.ClientSecret != nil {
		err := p.dataValidator.validateWebhookClientSecret(*data.ClientSecret, 12)
		if err != nil {
			return err
		}
	}

	if data.EndpointUrl != nil {
		err := p.dataValidator.validateWebhookEndpointUrl(*data.EndpointUrl)
		if err != nil {
			return err
		}
	}
	return nil
}
