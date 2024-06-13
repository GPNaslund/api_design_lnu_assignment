package dogsheltersservice

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type PostDogSheltersRepository interface {
	CreateDogShelter(ctx context.Context, newDogShelter dogshelterdto.NewDogShelterDTO) (model.DogShelter, error)
	GetDogShelterByUsername(ctx context.Context, username string) (model.DogShelter, error)
}

type PostDogSheltersCryptographyService interface {
	HashPassword(unhashedPassword string) (string, error)
}

type PostDogSheltersLinkGenerator interface {
	GenerateDogsFromDogShelterLink(shelterId string) string
	GenerateShelterLink(shelterId string) string
	GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO
}

type PostDogSheltersService struct {
	repo          PostDogSheltersRepository
	linkGenerator PostDogSheltersLinkGenerator
	cryptoService PostDogSheltersCryptographyService
}

func NewPostDogSheltersService(repo PostDogSheltersRepository,
	linkGenerator PostDogSheltersLinkGenerator, cryptoService PostDogSheltersCryptographyService) PostDogSheltersService {
	return PostDogSheltersService{
		repo:          repo,
		linkGenerator: linkGenerator,
		cryptoService: cryptoService,
	}
}

func (p PostDogSheltersService) CreateDogShelter(ctx context.Context,
	newDogShelter dogshelterdto.NewDogShelterDTO, credentials dto.UserCredentials) (dogshelterdto.DogShelterDTO, error) {
	emptyDto := dogshelterdto.DogShelterDTO{}
	role := credentials.UserRole
	if role == model.ADMIN {
		err := p.validateNewDogShelterDto(newDogShelter)
		if err != nil {
			return emptyDto, err
		}

		_, err = p.repo.GetDogShelterByUsername(ctx, *newDogShelter.Username)
		if err != nil {
			var dogShelterNotFoundError *customerrors.DogShelterNotFoundError
			if !errors.As(err, &dogShelterNotFoundError) {
				return emptyDto, err
			}
		} else {
			return emptyDto, &customerrors.InvalidNewDogShelterDataError{}
		}

		hashedPassword, err := p.cryptoService.HashPassword(*newDogShelter.Password)
		if err != nil {
			return emptyDto, &customerrors.CryptographyError{}
		}
		newDogShelter.Password = &hashedPassword

		dogShelter, err := p.repo.CreateDogShelter(ctx, newDogShelter)
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

	return emptyDto, &customerrors.UnauthorizedError{Message: "user role not recognized"}
}

func (p PostDogSheltersService) validateNewDogShelterDto(newDogShelter dogshelterdto.NewDogShelterDTO) error {
	var errMsgs []string
	stringPointerFields := map[string]*string{
		"name":     newDogShelter.Name,
		"website":  newDogShelter.Website,
		"country":  newDogShelter.Country,
		"city":     newDogShelter.City,
		"address":  newDogShelter.Address,
		"username": newDogShelter.Username,
		"password": newDogShelter.Password,
	}
	for fieldName, fieldValue := range stringPointerFields {
		if fieldValue == nil || *fieldValue == "" {
			errMsgs = append(errMsgs, fmt.Sprintf("%s cannot be empty", fieldName))
		}
	}

	if len(errMsgs) > 0 {
		return &customerrors.IncompleteDogShelterDataError{Message: strings.Join(errMsgs, " + ")}
	}

	return nil
}
