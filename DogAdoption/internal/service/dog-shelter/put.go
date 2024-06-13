package dogsheltersservice

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type PutDogSheltersRepository interface {
	UpdateDogShelter(ctx context.Context, dogId int, updateDogData dogshelterdto.UpdateDogShelterDTO) (model.DogShelter, error)
	GetDogShelterById(ctx context.Context, dogId int) (model.DogShelter, error)
}

type PutDogSheltersLinkGenerator interface {
	GenerateDogsFromDogShelterLink(shelterId string) string
	GenerateShelterLink(shelterId string) string
	GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO
}

type PutDogSheltersService struct {
	repo          PutDogSheltersRepository
	linkGenerator PutDogSheltersLinkGenerator
}

func NewPutDogSheltersService(repo PutDogSheltersRepository, linkGenerator PutDogSheltersLinkGenerator) PutDogSheltersService {
	return PutDogSheltersService{
		repo:          repo,
		linkGenerator: linkGenerator,
	}
}

func (p PutDogSheltersService) UpdateDogShelter(ctx context.Context,
	dogShelterId string, credentials dto.UserCredentials, updateDogShelter dogshelterdto.UpdateDogShelterDTO) (dogshelterdto.DogShelterDTO, error) {
	emptyDto := dogshelterdto.DogShelterDTO{}
	dogShelterIdInt, err := strconv.Atoi(dogShelterId)
	if err != nil {
		return emptyDto, &customerrors.IntegerConversionError{Message: "dog shelter id parameter needs to be a number"}
	}

	role := credentials.UserRole

	if role != model.ADMIN && role != model.DOGSHELTER {
		return emptyDto, &customerrors.UnauthorizedError{Message: "user role not recognized"}
	}

	if role == model.DOGSHELTER {
		dogShelter, err := p.repo.GetDogShelterById(ctx, dogShelterIdInt)
		if err != nil {
			return emptyDto, err
		}
		if dogShelter.Id != credentials.Id {
			return emptyDto, &customerrors.UnauthorizedError{Message: "unauthorized to update dog"}
		}
	}
	err = p.validateUpdateDogShelterDto(updateDogShelter)
	if err != nil {
		return emptyDto, err
	}

	dogShelter, err := p.repo.UpdateDogShelter(ctx, dogShelterIdInt, updateDogShelter)
	if err != nil {
		return emptyDto, err
	}
	var dogShelterDto dogshelterdto.DogShelterDTO
	dogShelterJson, err := json.Marshal(dogShelter.ToJson())
	if err != nil {
		return emptyDto, err
	}
	err = json.Unmarshal(dogShelterJson, &dogShelterDto)
	if err != nil {
		return emptyDto, err
	}

	selfLink := p.linkGenerator.GenerateShelterLink(fmt.Sprintf("%d", dogShelter.Id))
	dogsLink := p.linkGenerator.GenerateDogsFromDogShelterLink(fmt.Sprintf("%d", dogShelter.Id))
	dogShelterDto.Links = dogshelterdto.DogShelterDtoLinks{
		SelfLink: selfLink,
		DogsLink: dogsLink,
	}
	return dogShelterDto, nil
}

func (p PutDogSheltersService) validateUpdateDogShelterDto(updateDto dogshelterdto.UpdateDogShelterDTO) error {
	if updateDto.Name == nil && updateDto.Website == nil && updateDto.Country == nil && updateDto.City == nil && updateDto.Address == nil {
		return &customerrors.IncompleteDogShelterDataError{}
	}
	var errMsgs []string
	if updateDto.Name != nil {
		if *updateDto.Name == "" {
			errMsgs = append(errMsgs, "name cannot be empty")
		}
	}
	if updateDto.Website != nil {
		if *updateDto.Website == "" {
			errMsgs = append(errMsgs, "website cannot be empty")
		}
	}
	if updateDto.Country != nil {
		if *updateDto.Country == "" {
			errMsgs = append(errMsgs, "country cannot be empty")
		}
	}
	if updateDto.City != nil {
		if *updateDto.City == "" {
			errMsgs = append(errMsgs, "city cannot be empty")
		}
	}
	if updateDto.Address != nil {
		if *updateDto.Address == "" {
			errMsgs = append(errMsgs, "address cannot be empty")
		}
	}

	if len(errMsgs) > 0 {
		return &customerrors.IncompleteDogShelterDataError{Message: strings.Join(errMsgs, " + ")}
	}

	return nil
}
