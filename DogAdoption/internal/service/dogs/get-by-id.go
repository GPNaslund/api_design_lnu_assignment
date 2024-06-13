package dogsservice

import (
	dogdto "1dv027/aad/internal/dto/dog"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

type GetDogByIdRepository interface {
	GetDogById(ctx context.Context, dogId int) (model.Dog, error)
}

type GetDogByIdLinkGenerator interface {
	GenerateDogLink(dogId string) string
	GenerateShelterLink(shelterId string) string
}

type GetDogByIdService struct {
	repo          GetDogByIdRepository
	linkGenerator GetDogByIdLinkGenerator
}

func NewGetDogByIdService(repo GetDogByIdRepository, linkGenerator GetDogByIdLinkGenerator) GetDogByIdService {
	return GetDogByIdService{
		repo:          repo,
		linkGenerator: linkGenerator,
	}
}

func (g GetDogByIdService) GetDogById(ctx context.Context, dogId string) (dogdto.DogDTO, error) {
	emptyDto := dogdto.DogDTO{}
	dogIdInt, err := strconv.Atoi(dogId)
	if err != nil {
		return emptyDto, &customerrors.IntegerConversionError{}
	}

	dog, err := g.repo.GetDogById(ctx, dogIdInt)
	if err != nil {
		return emptyDto, err
	}

	dogJson, err := json.Marshal(dog.ToJson())
	if err != nil {
		return emptyDto, err
	}

	var dogDto dogdto.DogDTO
	err = json.Unmarshal(dogJson, &dogDto)
	if err != nil {
		return emptyDto, err
	}
	dogDto.Links = dogdto.DogLinksDTO{
		ShelterLink: g.linkGenerator.GenerateShelterLink(fmt.Sprintf("%d", dogDto.ShelterId)),
		SelfLink:    g.linkGenerator.GenerateDogLink(fmt.Sprintf("%d", dogDto.Id)),
	}

	return dogDto, nil
}
