package dogsheltersservice

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	"context"
	"encoding/json"
	"fmt"
)

type GetDogSheltersRepository interface {
	GetDogShelters(ctx context.Context, queryParams dto.QueryParams) (dogshelterdto.GetDogSheltersQueryResponseDTO, error)
}

type GetDogSheltersLinkGenerator interface {
	GenerateDogsFromDogShelterLink(shelterId string) string
	GenerateShelterLink(shelterId string) string
	GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO
}

type GetDogSheltersService struct {
	repo          GetDogSheltersRepository
	linkGenerator GetDogSheltersLinkGenerator
}

func NewGetDogSheltersService(repo GetDogSheltersRepository, linkGenerator GetDogSheltersLinkGenerator) GetDogSheltersService {
	return GetDogSheltersService{
		repo:          repo,
		linkGenerator: linkGenerator,
	}
}

func (g GetDogSheltersService) GetShelters(ctx context.Context,
	queryParams dto.QueryParams) (dogshelterdto.DogSheltersAndPaginationLinksDTO, error) {
	emptyDto := dogshelterdto.DogSheltersAndPaginationLinksDTO{}
	shelterResult, err := g.repo.GetDogShelters(ctx, queryParams)
	if err != nil {
		return emptyDto, err
	}
	var shelterDtoSlice []dogshelterdto.DogShelterDTO
	for _, shelter := range shelterResult.DogShelters {
		var shelterDto dogshelterdto.DogShelterDTO
		shelterJson, err := json.Marshal(shelter.ToJson())
		if err != nil {
			return emptyDto, err
		}
		err = json.Unmarshal(shelterJson, &shelterDto)
		if err != nil {
			return emptyDto, err
		}

		selfLink := g.linkGenerator.GenerateShelterLink(fmt.Sprintf("%d", shelterDto.Id))
		dogsLink := g.linkGenerator.GenerateDogsFromDogShelterLink(fmt.Sprintf("%d", shelterDto.Id))
		shelterDto.Links = dogshelterdto.DogShelterDtoLinks{
			SelfLink: selfLink,
			DogsLink: dogsLink,
		}
		shelterDtoSlice = append(shelterDtoSlice, shelterDto)
	}
	paginationLinks := g.linkGenerator.GeneratePaginationLinks(shelterResult.TotalAmountAvailable, queryParams, "/dogshelters")
	dogShelterAndPagination := dogshelterdto.DogSheltersAndPaginationLinksDTO{
		DogShelterData:  shelterDtoSlice,
		PaginationLinks: paginationLinks,
	}
	return dogShelterAndPagination, nil
}
