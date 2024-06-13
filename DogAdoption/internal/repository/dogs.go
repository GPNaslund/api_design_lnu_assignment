package repository

import (
	"1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	"1dv027/aad/internal/model"
	"context"
)

type DogsDataAccess interface {
	GetDogs(ctx context.Context, queryParams dto.QueryParams) (dogdto.GetDogsQueryResponseDTO, error)
	GetDogById(ctx context.Context, dogId int) (model.Dog, error)
	DeleteDog(ctx context.Context, dogId int) error
	UpdateDog(ctx context.Context, dogId int, updatedDogData dogdto.UpdateDogDTO) error
	CreateDog(ctx context.Context, newDog dogdto.NewDogDTO) (int, error)
}

type DogsRepository struct {
	dogsDataAccess DogsDataAccess
}

func NewDogsRepository(dogsDataAccess DogsDataAccess) *DogsRepository {
	return &DogsRepository{
		dogsDataAccess: dogsDataAccess,
	}
}

func (d DogsRepository) GetDogs(ctx context.Context, params dto.QueryParams) (dogdto.GetDogsQueryResponseDTO, error) {
	return d.dogsDataAccess.GetDogs(ctx, params)
}

func (d DogsRepository) GetDogById(ctx context.Context, dogId int) (model.Dog, error) {
	return d.dogsDataAccess.GetDogById(ctx, dogId)

}

func (d DogsRepository) DeleteDog(ctx context.Context, dogId int) error {
	return d.dogsDataAccess.DeleteDog(ctx, dogId)
}

func (d DogsRepository) UpdateDog(ctx context.Context, dogId int, updatedDogData dogdto.UpdateDogDTO) (model.Dog, error) {
	err := d.dogsDataAccess.UpdateDog(ctx, dogId, updatedDogData)
	if err != nil {
		return model.Dog{}, err
	}
	return d.dogsDataAccess.GetDogById(ctx, dogId)
}

func (d DogsRepository) CreateDog(ctx context.Context, newDog dogdto.NewDogDTO) (model.Dog, error) {
	id, err := d.dogsDataAccess.CreateDog(ctx, newDog)
	if err != nil {
		return model.Dog{}, err
	}
	dog, err := d.dogsDataAccess.GetDogById(ctx, id)
	if err != nil {
		return model.Dog{}, err
	}
	return dog, nil
}
