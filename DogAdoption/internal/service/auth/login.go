package authservice

import (
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"errors"
)

type LoginRepository interface {
	GetAdminByUsername(ctx context.Context, username string) (model.Admin, error)
	GetDogShelterByUsername(ctx context.Context, username string) (model.DogShelter, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}

type JwtGenerator interface {
	GenerateJwt(username string, id int, userType model.UserRole) (string, error)
}

type CryptographyService interface {
	ComparePasswords(hashedPassword, passwordAttempt string) error
}

type LoginService struct {
	loginRepo     LoginRepository
	jwtGenerator  JwtGenerator
	cryptoService CryptographyService
}

func NewLoginService(loginRepo LoginRepository, jwtGenerator JwtGenerator, cryptoService CryptographyService) *LoginService {
	return &LoginService{
		loginRepo:     loginRepo,
		jwtGenerator:  jwtGenerator,
		cryptoService: cryptoService,
	}
}

func (l LoginService) ValidateUsernameAndPassword(ctx context.Context, username, password string) (string, error) {

	admin, err := l.loginRepo.GetAdminByUsername(ctx, username)
	if err != nil {
		var adminNotFoundError *customerrors.AdminNotFoundError
		if !errors.As(err, &adminNotFoundError) {
			return "", err
		}
	} else {
		err := l.cryptoService.ComparePasswords(admin.Password, password)
		if err != nil {
			return "", &customerrors.WrongCredentialsError{}
		}
		return l.generateJwtForRole(admin.Username, admin.Id, model.ADMIN)
	}

	dogShelter, err := l.loginRepo.GetDogShelterByUsername(ctx, username)
	if err != nil {
		var dogShelterNotFoundError *customerrors.DogShelterNotFoundError
		if !errors.As(err, &dogShelterNotFoundError) {
			return "", err
		}
	} else {
		err := l.cryptoService.ComparePasswords(dogShelter.Password, password)
		if err != nil {
			return "", &customerrors.WrongCredentialsError{}
		}
		return l.generateJwtForRole(dogShelter.Username, dogShelter.Id, model.DOGSHELTER)
	}

	user, err := l.loginRepo.GetUserByUsername(ctx, username)
	if err != nil {
		var userNotFound *customerrors.UserNotFoundError
		if !errors.As(err, &userNotFound) {
			return "", err
		}
	} else {
		err := l.cryptoService.ComparePasswords(user.Password, password)
		if err != nil {
			return "", &customerrors.WrongCredentialsError{}
		}
		return l.generateJwtForRole(user.Username, user.Id, model.USER)
	}

	return "", &customerrors.UnauthorizedError{}
}

func (l LoginService) GetAllowedFields() map[string]any {
	return map[string]any{
		"username": "",
		"password": "",
	}
}

func (l LoginService) generateJwtForRole(username string, userId int, role model.UserRole) (string, error) {
	jwt, err := l.jwtGenerator.GenerateJwt(username, userId, role)
	if err != nil {
		return "", &customerrors.JwtError{}
	}
	return jwt, nil
}
