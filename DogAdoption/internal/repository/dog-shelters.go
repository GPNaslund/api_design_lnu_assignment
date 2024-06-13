package repository

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	"1dv027/aad/internal/model"
	"context"
)

type DogSheltersDataAccess interface {
	GetDogShelters(ctx context.Context, queryParams dto.QueryParams) (dogshelterdto.GetDogSheltersQueryResponseDTO, error)
	GetDogShelterById(ctx context.Context, shelterId int) (model.DogShelter, error)
	GetDogShelterByUsername(ctx context.Context, username string) (model.DogShelter, error)
	DeleteDogShelter(ctx context.Context, shelterId int) error
	UpdateDogShelter(ctx context.Context, shelterId int, updatedDogShelterData dogshelterdto.UpdateDogShelterDTO) error
	CreateDogShelter(ctx context.Context, newShelter dogshelterdto.NewDogShelterDTO) (int, error)
}

type DogSheltersRepository struct {
	dataAccess DogSheltersDataAccess
}

func NewDogSheltersRepository(dataAccess DogSheltersDataAccess) DogSheltersRepository {
	return DogSheltersRepository{
		dataAccess: dataAccess,
	}
}

func (d DogSheltersRepository) GetDogShelters(ctx context.Context, params dto.QueryParams) (dogshelterdto.GetDogSheltersQueryResponseDTO, error) {
	return d.dataAccess.GetDogShelters(ctx, params)
}

func (d DogSheltersRepository) GetDogShelterById(ctx context.Context, dogShelterId int) (model.DogShelter, error) {
	return d.dataAccess.GetDogShelterById(ctx, dogShelterId)
}

func (d DogSheltersRepository) GetDogShelterByUsername(ctx context.Context, username string) (model.DogShelter, error) {
	return d.dataAccess.GetDogShelterByUsername(ctx, username)
}

func (d DogSheltersRepository) DeleteDogShelter(ctx context.Context, dogShelterId int) error {
	return d.dataAccess.DeleteDogShelter(ctx, dogShelterId)
}

func (d DogSheltersRepository) UpdateDogShelter(ctx context.Context, dogShelterId int, dogShelter dogshelterdto.UpdateDogShelterDTO) (model.DogShelter, error) {
	err := d.dataAccess.UpdateDogShelter(ctx, dogShelterId, dogShelter)
	if err != nil {
		return model.DogShelter{}, err
	}
	return d.dataAccess.GetDogShelterById(ctx, dogShelterId)
}

func (d DogSheltersRepository) CreateDogShelter(ctx context.Context, dogShelter dogshelterdto.NewDogShelterDTO) (model.DogShelter, error) {
	emptyDto := model.DogShelter{}
	id, err := d.dataAccess.CreateDogShelter(ctx, dogShelter)
	if err != nil {
		return emptyDto, err
	}
	dog, err := d.dataAccess.GetDogShelterById(ctx, id)
	if err != nil {
		return emptyDto, err
	}
	return dog, nil
}
