package dogsservice

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"strconv"
)

type DeleteDogRepository interface {
	DeleteDog(ctx context.Context, dogId int) error
	GetDogById(ctx context.Context, dogId int) (model.Dog, error)
}

type DeleteDogService struct {
	repo DeleteDogRepository
}

func NewDeleteDogService(repo DeleteDogRepository) DeleteDogService {
	return DeleteDogService{
		repo: repo,
	}
}

func (d DeleteDogService) DeleteDog(ctx context.Context, dogId string, credentials dto.UserCredentials) error {
	dogIdInt, err := strconv.Atoi(dogId)
	if err != nil {
		return &customerrors.IntegerConversionError{}
	}
	role := credentials.UserRole
	if role == model.ADMIN {
		return d.repo.DeleteDog(ctx, dogIdInt)
	} else if role == model.DOGSHELTER {
		dog, err := d.repo.GetDogById(ctx, dogIdInt)
		if err != nil {
			return err
		}
		if dog.ShelterId != credentials.Id {
			return &customerrors.UnauthorizedError{}
		}
		return d.repo.DeleteDog(ctx, dogIdInt)
	}
	return &customerrors.UnauthorizedError{}
}
