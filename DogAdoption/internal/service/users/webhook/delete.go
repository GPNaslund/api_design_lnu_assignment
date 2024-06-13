package userwebhookservice

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"strconv"
)

type DeleteUserWebhookRepository interface {
	DeleteUserWebhook(ctx context.Context, userId int) error
}

type DeleteUserWebhookService struct {
	repo DeleteUserWebhookRepository
}

func NewDeleteWebhookService(repo DeleteUserWebhookRepository) DeleteUserWebhookService {
	return DeleteUserWebhookService{
		repo: repo,
	}
}

func (d DeleteUserWebhookService) DeleteWebhook(ctx context.Context, idParam string, user dto.UserCredentials) error {
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		return &customerrors.IntegerConversionError{}
	}

	if user.UserRole == model.ADMIN {
		err := d.repo.DeleteUserWebhook(ctx, idParamInt)
		if err != nil {
			return err
		}
		return nil
	}

	if user.UserRole == model.USER {
		if idParamInt != user.Id {
			return &customerrors.UnauthorizedError{}
		}
		err := d.repo.DeleteUserWebhook(ctx, idParamInt)
		if err != nil {
			return err
		}
		return nil
	}

	return &customerrors.UnauthorizedError{}
}
