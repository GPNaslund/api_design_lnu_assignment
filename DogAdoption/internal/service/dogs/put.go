package dogsservice

import (
	"1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

type PutDogsRepository interface {
	UpdateDog(ctx context.Context, dogId int, updateDogData dogdto.UpdateDogDTO) (model.Dog, error)
	GetDogById(ctx context.Context, dogId int) (model.Dog, error)
}

type PutDogsLinkGenerator interface {
	GenerateDogLink(dogId string) string
	GenerateShelterLink(shelterId string) string
	GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO
}

type PutDogService struct {
	repo          PutDogsRepository
	linkGenerator PutDogsLinkGenerator
}

func NewPutDogService(repo PutDogsRepository, linkGenerator PutDogsLinkGenerator) PutDogService {
	return PutDogService{
		repo:          repo,
		linkGenerator: linkGenerator,
	}
}

func (p PutDogService) UpdateDog(ctx context.Context,
	dogId string, credentials dto.UserCredentials, updatedDog dogdto.UpdateDogDTO) (dogdto.DogDTO, error) {
	emptyDto := dogdto.DogDTO{}
	dogIdInt, err := strconv.Atoi(dogId)
	if err != nil {
		return emptyDto, &customerrors.IntegerConversionError{}
	}

	role := credentials.UserRole

	if role != model.ADMIN && role != model.DOGSHELTER {
		return emptyDto, &customerrors.UnauthorizedError{}
	}

	if role == model.DOGSHELTER {
		dog, err := p.repo.GetDogById(ctx, dogIdInt)
		if err != nil {
			return emptyDto, err
		}
		if dog.ShelterId != credentials.Id {
			return emptyDto, &customerrors.UnauthorizedError{}
		}
	}

	err = p.validateGenderField(updatedDog)
	if err != nil {
		return emptyDto, err
	}

	dogModel, err := p.repo.UpdateDog(ctx, dogIdInt, updatedDog)
	if err != nil {
		return emptyDto, err
	}
	dogJson, err := json.Marshal(dogModel)
	if err != nil {
		return emptyDto, err
	}

	var dogDto dogdto.DogDTO
	err = json.Unmarshal(dogJson, &dogDto)
	if err != nil {
		return emptyDto, err
	}
	dogDto.Links = dogdto.DogLinksDTO{
		ShelterLink: p.linkGenerator.GenerateShelterLink(fmt.Sprintf("%d", dogDto.ShelterId)),
		SelfLink:    p.linkGenerator.GenerateDogLink(fmt.Sprintf("%d", dogDto.Id)),
	}

	return dogDto, nil
}

func (p PutDogService) validateGenderField(updatedDog dogdto.UpdateDogDTO) error {
	if updatedDog.Gender != nil {
		if *updatedDog.Gender != "male" && *updatedDog.Gender != "female" {
			return &customerrors.IncompleteDogDataError{}
		}
		return nil
	}
	return nil
}
