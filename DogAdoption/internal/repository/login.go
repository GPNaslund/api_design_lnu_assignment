package repository

import (
	"1dv027/aad/internal/model"
	"context"
)

type GetAdminsDataAccess interface {
	GetAdminByUsername(ctx context.Context, username string) (model.Admin, error)
}

type GetDogSheltersDataAccess interface {
	GetDogShelterByUsername(ctx context.Context, username string) (model.DogShelter, error)
}

type GetUsersDataAccess interface {
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}

type LoginRepository struct {
	adminsDataAccess      GetAdminsDataAccess
	dogSheltersDataAccess GetDogSheltersDataAccess
	usersDataAccess       GetUsersDataAccess
}

func NewLoginRepository(adminsDataAccess GetAdminsDataAccess,
	dogSheltersDataAccess GetDogSheltersDataAccess, usersDataAccess GetUsersDataAccess) *LoginRepository {
	return &LoginRepository{
		adminsDataAccess:      adminsDataAccess,
		dogSheltersDataAccess: dogSheltersDataAccess,
		usersDataAccess:       usersDataAccess,
	}
}

func (l LoginRepository) GetDogShelterByUsername(ctx context.Context, username string) (model.DogShelter, error) {
	emptyModel := model.DogShelter{}
	dogShelter, err := l.dogSheltersDataAccess.GetDogShelterByUsername(ctx, username)
	if err != nil {
		return emptyModel, err
	}
	return dogShelter, nil
}

func (l LoginRepository) GetAdminByUsername(ctx context.Context, username string) (model.Admin, error) {
	emptyModel := model.Admin{}
	admin, err := l.adminsDataAccess.GetAdminByUsername(ctx, username)
	if err != nil {
		return emptyModel, err
	}
	return admin, nil
}

func (l LoginRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	emptyModel := model.User{}
	user, err := l.usersDataAccess.GetUserByUsername(ctx, username)
	if err != nil {
		return emptyModel, err
	}
	return user, nil
}
