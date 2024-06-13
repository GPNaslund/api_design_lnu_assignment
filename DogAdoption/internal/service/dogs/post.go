package dogsservice

import (
	"1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	customerrors "1dv027/aad/internal/errors"
	"1dv027/aad/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type PostDogsRepository interface {
	CreateDog(ctx context.Context, newDog dogdto.NewDogDTO) (model.Dog, error)
}

type PostDogsLinkGenerator interface {
	GenerateDogLink(dogId string) string
	GenerateShelterLink(shelterId string) string
	GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO
}

type NewDogWebhookDispatcher interface {
	DispatchNewDogWebhook(ctx context.Context, dogData dogdto.DogDTO)
}

type PostDogService struct {
	repo              PostDogsRepository
	linkGenerator     PostDogsLinkGenerator
	webhookDispatcher NewDogWebhookDispatcher
}

func NewPostDogService(repo PostDogsRepository, linkGenerator PostDogsLinkGenerator, webhookDispatcher NewDogWebhookDispatcher) PostDogService {
	return PostDogService{
		repo:              repo,
		linkGenerator:     linkGenerator,
		webhookDispatcher: webhookDispatcher,
	}
}

func (p PostDogService) CreateDog(ctx context.Context,
	newDog dogdto.NewDogDTO, credentials dto.UserCredentials) (dogdto.DogDTO, error) {
	emptyDto := dogdto.DogDTO{}
	role := credentials.UserRole

	if role != model.ADMIN && role != model.DOGSHELTER {
		return emptyDto, &customerrors.UnauthorizedError{Message: "user role not recognized"}
	}

	err := p.validateFields(newDog)
	if err != nil {
		return emptyDto, err
	}

	if role == model.DOGSHELTER {
		newDog.ShelterId = &credentials.Id
	}

	dog, err := p.repo.CreateDog(ctx, newDog)
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
		ShelterLink: p.linkGenerator.GenerateShelterLink(fmt.Sprintf("%d", dogDto.ShelterId)),
		SelfLink:    p.linkGenerator.GenerateDogLink(fmt.Sprintf("%d", dogDto.Id)),
	}

	p.webhookDispatcher.DispatchNewDogWebhook(ctx, dogDto)
	return dogDto, nil
}

func (p PostDogService) validateFields(newDog dogdto.NewDogDTO) error {
	var errMsgs []string
	stringPointerFields := map[string]*string{
		"name":          newDog.Name,
		"description":   newDog.Description,
		"breed":         newDog.Breed,
		"image_url":     newDog.ImageUrl,
		"friendly_with": newDog.FriendlyWith,
		"gender":        newDog.Gender,
	}
	for fieldName, fieldValue := range stringPointerFields {
		if fieldValue == nil || *fieldValue == "" {
			errMsgs = append(errMsgs, fmt.Sprintf("%s cannot be empty", fieldName))
		}
	}

	intPointerFields := map[string]*int{
		"shelter_id":   newDog.ShelterId,
		"adoption_fee": newDog.AdoptionFee,
	}
	for fieldName, fieldValue := range intPointerFields {
		if fieldValue == nil {
			errMsgs = append(errMsgs, fmt.Sprintf("%s cannot be empty", fieldName))
		}
	}

	boolPointerFields := map[string]*bool{
		"is_neutered": newDog.IsNeutered,
		"is_adopted":  newDog.IsAdopted,
	}
	for fieldName, fieldValue := range boolPointerFields {
		if fieldValue == nil {
			errMsgs = append(errMsgs, fmt.Sprintf("%s cannot be empty", fieldName))
		}
	}

	if newDog.Gender != nil && *newDog.Gender != "male" && *newDog.Gender != "female" {
		errMsgs = append(errMsgs, "gender must be either 'male' or 'female'")
	}

	if newDog.BirthDate == nil {
		errMsgs = append(errMsgs, "birth_date cannot be empty")
	} else {
		if newDog.BirthDate.After(time.Now()) {
			errMsgs = append(errMsgs, "birth_date must be in the past")
		}
	}

	if len(errMsgs) > 0 {
		return &customerrors.IncompleteDogDataError{Message: strings.Join(errMsgs, " + ")}
	}

	return nil
}
