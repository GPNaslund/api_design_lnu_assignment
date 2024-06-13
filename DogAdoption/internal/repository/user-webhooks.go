package repository

import (
	webhookdto "1dv027/aad/internal/dto/user/webhook"
	"1dv027/aad/internal/model"
	"context"
)

type UserWebhooksDataAccess interface {
	DeleteUserWebhook(ctx context.Context, userId int) error
	GetUserWebhook(ctx context.Context, userId int) (model.Webhook, error)
	CreateNewWebhook(ctx context.Context, userId int, data webhookdto.NewUserWebhookDTO) error
	UpdateUserWebhook(ctx context.Context, userId int, data webhookdto.UpdateUserWebhookDTO) error
	GetAllWebhooksByAction(ctx context.Context, action model.WebhookAction) ([]model.Webhook, error)
}

type UserWebhooksRepository struct {
	dataaccess UserWebhooksDataAccess
}

func NewUserWebhooksRepository(dataaccess UserWebhooksDataAccess) UserWebhooksRepository {
	return UserWebhooksRepository{
		dataaccess: dataaccess,
	}
}

func (u UserWebhooksRepository) DeleteUserWebhook(ctx context.Context, userId int) error {
	return u.dataaccess.DeleteUserWebhook(ctx, userId)
}

func (u UserWebhooksRepository) GetUserWebhook(ctx context.Context, userId int) (model.Webhook, error) {
	return u.dataaccess.GetUserWebhook(ctx, userId)
}

func (u UserWebhooksRepository) CreateNewUserWebhook(ctx context.Context, userId int, data webhookdto.NewUserWebhookDTO) (model.Webhook, error) {
	emptyModel := model.Webhook{}
	err := u.dataaccess.CreateNewWebhook(ctx, userId, data)
	if err != nil {
		return emptyModel, err
	}
	webhook, err := u.dataaccess.GetUserWebhook(ctx, userId)
	if err != nil {
		return emptyModel, err
	}
	return webhook, nil
}

func (u UserWebhooksRepository) UpdateUserWebhook(ctx context.Context, userId int, data webhookdto.UpdateUserWebhookDTO) (model.Webhook, error) {
	emptyModel := model.Webhook{}
	err := u.dataaccess.UpdateUserWebhook(ctx, userId, data)
	if err != nil {
		return emptyModel, err
	}
	webhookModel, err := u.dataaccess.GetUserWebhook(ctx, userId)
	if err != nil {
		return emptyModel, err
	}
	return webhookModel, nil
}

func (u UserWebhooksRepository) GetAllWebhooksByAction(ctx context.Context, action model.WebhookAction) ([]model.Webhook, error) {
	return u.dataaccess.GetAllWebhooksByAction(ctx, action)
}
