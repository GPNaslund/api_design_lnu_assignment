package usersservice

import (
	userdto "1dv027/aad/internal/dto/user"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
)

type PostUsersRepository interface {
	CreateNewUser(ctx context.Context, newUser userdto.NewUserDTO) (model.User, error)
}

type PostUsersCryptographyService interface {
	HashPassword(unhashedPassword string) (string, error)
}

type PostUsersService struct {
	userRepo      PostUsersRepository
	cryptoService PostUsersCryptographyService
}

func NewPostUsersService(userRepo PostUsersRepository, cryptoService PostUsersCryptographyService) PostUsersService {
	return PostUsersService{
		userRepo:      userRepo,
		cryptoService: cryptoService,
	}
}

func (p PostUsersService) CreateNewUser(ctx context.Context, newUser userdto.NewUserDTO) (userdto.UserDTO, error) {
	emptyDto := userdto.UserDTO{}
	err := p.validateNewUserDto(newUser)
	if err != nil {
		return emptyDto, err
	}

	hashedPassword, err := p.cryptoService.HashPassword(*newUser.Password)
	if err != nil {
		return emptyDto, &customerrors.CryptographyError{}
	}
	newUser.Password = &hashedPassword
	user, err := p.userRepo.CreateNewUser(ctx, newUser)
	if err != nil {
		return emptyDto, err
	}

	var userDto userdto.UserDTO
	userJson, err := json.Marshal(user.ToJson())
	if err != nil {
		return emptyDto, err
	}
	err = json.Unmarshal(userJson, &userDto)
	if err != nil {
		return emptyDto, err
	}

	return userDto, nil
}

func (p PostUsersService) validateNewUserDto(dto userdto.NewUserDTO) error {
	if dto.Username == nil || *dto.Username == "" {
		return &customerrors.IncompleteNewUserError{}
	}
	if dto.Password == nil || *dto.Password == "" {
		return &customerrors.IncompleteNewUserError{}
	}

	return nil
}
