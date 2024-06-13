package repository

import (
	"1dv027/aad/internal/dto"
	userdto "1dv027/aad/internal/dto/user"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"errors"
)

type UsersDataAccess interface {
	CreateNewUser(ctx context.Context, newUser userdto.NewUserDTO) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	DeleteUser(ctx context.Context, userId int) error
}

type UsersRepository struct {
	usersDataAccess UsersDataAccess
}

func NewUsersRepository(usersDataAccess UsersDataAccess) UsersRepository {
	return UsersRepository{
		usersDataAccess: usersDataAccess,
	}
}

func (u UsersRepository) CreateNewUser(ctx context.Context, newUser userdto.NewUserDTO) (model.User, error) {
	existingUser, err := u.usersDataAccess.GetUserByUsername(ctx, *newUser.Username)

	if err != nil {
		var userNotFoundError *customerrors.UserNotFoundError
		if !errors.As(err, &userNotFoundError) {
			return model.User{}, err
		}
	} else if existingUser.Username != "" {
		return model.User{}, &customerrors.InvalidNewUserDataError{Message: "invalid username. try another one!"}
	}

	createdUser, err := u.usersDataAccess.CreateNewUser(ctx, newUser)
	if err != nil {
		return model.User{}, err
	}
	return createdUser, nil
}

func (u UsersRepository) DeleteUser(ctx context.Context, userId int) error {
	return u.usersDataAccess.DeleteUser(ctx, userId)
}

func (u UsersRepository) GetAuthenticatedUser(ctx context.Context, user dto.UserCredentials) (model.User, error) {
	emptyModel := model.User{}
	userModel, err := u.usersDataAccess.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return emptyModel, err
	}
	return userModel, nil
}
