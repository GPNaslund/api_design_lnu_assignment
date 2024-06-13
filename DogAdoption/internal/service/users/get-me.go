package usersservice

import (
	"1dv027/aad/internal/dto"
	userdto "1dv027/aad/internal/dto/user"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
)

type GetUsersMeRepository interface {
	GetAuthenticatedUser(ctx context.Context, user dto.UserCredentials) (model.User, error)
}

type GetUsersMeService struct {
	repo GetUsersMeRepository
}

func NewGetUsersMeService(repo GetUsersMeRepository) GetUsersMeService {
	return GetUsersMeService{
		repo: repo,
	}
}

func (g GetUsersMeService) GetAuthenticatedUser(ctx context.Context, user dto.UserCredentials) (userdto.UserDTO, error) {
	emptyDto := userdto.UserDTO{}
	if user.UserRole != model.USER {
		return emptyDto, &customerrors.UnauthorizedError{}
	}
	userModel, err := g.repo.GetAuthenticatedUser(ctx, user)
	if err != nil {
		return emptyDto, err
	}

	userJson, err := json.Marshal(userModel.ToJson())
	if err != nil {
		return emptyDto, err
	}
	var userDto userdto.UserDTO
	err = json.Unmarshal(userJson, &userDto)
	if err != nil {
		return emptyDto, err
	}

	return userDto, nil
}
