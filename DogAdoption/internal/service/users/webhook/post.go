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

type PostUserWebhooksDataValidator interface {
	validateWebhookActions(actions []string) error
	validateWebhookEndpointUrl(urlString string) error
	validateWebhookClientSecret(secret string, minSecretLength int) error
}

type PostUserWebhooksRepository interface {
	CreateNewUserWebhook(ctx context.Context, userId int, data webhookdto.NewUserWebhookDTO) (model.Webhook, error)
}

type PostUserWebhooksCryptographyService interface {
	EncryptPlainText(plainText string) (string, error)
}

type PostUserWebhookService struct {
	repo          PostUserWebhooksRepository
	dataValidator PostUserWebhooksDataValidator
	cryptoService PostUserWebhooksCryptographyService
}

func NewPostUserWebhookService(repo PostUserWebhooksRepository,
	dataValidator PostUserWebhooksDataValidator, cryptoService PostUserWebhooksCryptographyService) PostUserWebhookService {
	return PostUserWebhookService{
		repo:          repo,
		dataValidator: dataValidator,
		cryptoService: cryptoService,
	}
}

func (p PostUserWebhookService) CreateNewWebhook(ctx context.Context, idParam string,
	user dto.UserCredentials, webhookData webhookdto.NewUserWebhookDTO) (webhookdto.UserWebhookDTO, error) {
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

	err = p.validateNewWebhookData(webhookData)
	if err != nil {
		return emptyDto, err
	}

	encryptedClientSecret, err := p.cryptoService.EncryptPlainText(*webhookData.ClientSecret)
	if err != nil {
		return emptyDto, &customerrors.CryptographyError{}
	}
	webhookData.ClientSecret = &encryptedClientSecret
	webhookModel, err := p.repo.CreateNewUserWebhook(ctx, idParamInt, webhookData)
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

func (p PostUserWebhookService) validateNewWebhookData(webhookData webhookdto.NewUserWebhookDTO) error {
	if webhookData.EndpointUrl == nil || webhookData.Actions == nil || webhookData.ClientSecret == nil {
		return &customerrors.IncompleteWebhookDataError{}
	}

	err := p.dataValidator.validateWebhookActions(*webhookData.Actions)
	if err != nil {
		return err
	}

	err = p.dataValidator.validateWebhookClientSecret(*webhookData.ClientSecret, 12)
	if err != nil {
		return err
	}

	err = p.dataValidator.validateWebhookEndpointUrl(*webhookData.EndpointUrl)
	if err != nil {
		return err
	}
	return nil
}
