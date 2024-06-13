package usersservice

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"strconv"
)

type DeleteUsersRepository interface {
	DeleteUser(ctx context.Context, userId int) error
}

type DeleteUsersService struct {
	repo DeleteUsersRepository
}

func NewDeleteUsersService(repo DeleteUsersRepository) DeleteUsersService {
	return DeleteUsersService{
		repo: repo,
	}
}

func (d DeleteUsersService) DeleteUser(ctx context.Context, idParam string, user dto.UserCredentials) error {
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		return &customerrors.IntegerConversionError{}
	}

	if user.UserRole == model.ADMIN {
		err := d.repo.DeleteUser(ctx, idParamInt)
		if err != nil {
			return err
		}
		return nil
	}

	if user.UserRole == model.USER && idParamInt == user.Id {
		err := d.repo.DeleteUser(ctx, idParamInt)
		if err != nil {
			return err
		}
		return nil
	}

	return &customerrors.UnauthorizedError{}
}
