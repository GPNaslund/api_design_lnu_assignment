package userwebhookservice

import (
	"1dv027/aad/internal/dto"
	userdto "1dv027/aad/internal/dto/user/webhook"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"strconv"
)

type GetUserWebhookRepository interface {
	GetUserWebhook(ctx context.Context, userId int) (model.Webhook, error)
}

type GetUserWebhookService struct {
	repo GetUserWebhookRepository
}

func NewGetUserWebhookService(repo GetUserWebhookRepository) GetUserWebhookService {
	return GetUserWebhookService{
		repo: repo,
	}
}

func (g GetUserWebhookService) GetUserWebhook(ctx context.Context, idParam string, userCredentials dto.UserCredentials) (userdto.UserWebhookDTO, error) {
	emptyDto := userdto.UserWebhookDTO{}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		return emptyDto, &customerrors.IntegerConversionError{}
	}

	if userCredentials.UserRole != model.ADMIN && userCredentials.UserRole != model.USER {
		return emptyDto, &customerrors.UnauthorizedError{}
	}

	if userCredentials.UserRole == model.USER && idParamInt != userCredentials.Id {
		return emptyDto, &customerrors.UnauthorizedError{}
	}

	webhookModel, err := g.repo.GetUserWebhook(ctx, idParamInt)
	if err != nil {
		return emptyDto, err
	}
	webhookJson, err := json.Marshal(webhookModel.ToJson())
	if err != nil {
		return emptyDto, err
	}
	var webhookDto userdto.UserWebhookDTO
	err = json.Unmarshal(webhookJson, &webhookDto)
	if err != nil {
		return emptyDto, err
	}
	return webhookDto, nil
}
