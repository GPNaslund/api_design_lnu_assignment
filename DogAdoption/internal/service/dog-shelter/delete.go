package dogsheltersservice

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"strconv"
)

type DeleteDogSheltersRepository interface {
	DeleteDogShelter(ctx context.Context, dogId int) error
	GetDogShelterById(ctx context.Context, dogId int) (model.DogShelter, error)
}

type DeleteDogSheltersService struct {
	repo DeleteDogSheltersRepository
}

func NewDeleteDogSheltersService(repo DeleteDogSheltersRepository) DeleteDogSheltersService {
	return DeleteDogSheltersService{
		repo: repo,
	}
}

func (d DeleteDogSheltersService) DeleteDogShelter(ctx context.Context, shelterId string, credentials dto.UserCredentials) error {
	shelterIdInt, err := strconv.Atoi(shelterId)
	if err != nil {
		return &customerrors.IntegerConversionError{}
	}
	role := credentials.UserRole
	if role == model.ADMIN {
		return d.repo.DeleteDogShelter(ctx, shelterIdInt)
	} else if role == model.DOGSHELTER {
		dogShelter, err := d.repo.GetDogShelterById(ctx, shelterIdInt)
		if err != nil {
			return err
		}
		if dogShelter.Id != credentials.Id {
			return &customerrors.UnauthorizedError{}
		}
		return d.repo.DeleteDogShelter(ctx, shelterIdInt)
	}
	return &customerrors.UnauthorizedError{}
}
