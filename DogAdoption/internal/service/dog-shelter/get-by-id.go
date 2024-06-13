package dogsheltersservice

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"strconv"
)

type GetDogSheltersByIdRepository interface {
	GetDogShelterById(ctx context.Context, dogId int) (model.DogShelter, error)
}

type GetDogSheltersByIdLinkGenerator interface {
	GenerateDogsFromDogShelterLink(shelterId string) string
	GenerateShelterLink(shelterId string) string
	GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO
}

type GetDogSheltersByIdService struct {
	repo          GetDogSheltersByIdRepository
	linkGenerator GetDogSheltersByIdLinkGenerator
}

func NewGetDogSheltersByIdService(repo GetDogSheltersByIdRepository, linkGenerator GetDogSheltersByIdLinkGenerator) GetDogSheltersByIdService {
	return GetDogSheltersByIdService{
		repo:          repo,
		linkGenerator: linkGenerator,
	}
}

func (g GetDogSheltersByIdService) GetDogShelterById(ctx context.Context, dogShelterId string) (dogshelterdto.DogShelterDTO, error) {
	emptyDto := dogshelterdto.DogShelterDTO{}
	dogShelterIdInt, err := strconv.Atoi(dogShelterId)
	if err != nil {
		return emptyDto, &customerrors.IntegerConversionError{}
	}
	dogShelter, err := g.repo.GetDogShelterById(ctx, dogShelterIdInt)
	if err != nil {
		return emptyDto, err
	}
	dogShelterJson, err := json.Marshal(dogShelter.ToJson())
	if err != nil {
		return emptyDto, err
	}
	var dogShelterDto dogshelterdto.DogShelterDTO
	err = json.Unmarshal(dogShelterJson, &dogShelterDto)
	if err != nil {
		return emptyDto, err
	}
	selfLink := g.linkGenerator.GenerateShelterLink(dogShelterId)
	dogsLink := g.linkGenerator.GenerateDogsFromDogShelterLink(dogShelterId)
	dogShelterDto.Links = dogshelterdto.DogShelterDtoLinks{
		SelfLink: selfLink,
		DogsLink: dogsLink,
	}
	return dogShelterDto, nil
}
