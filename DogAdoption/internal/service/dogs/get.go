package dogsservice

import (
	dto "1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	"context"
	"encoding/json"
	"fmt"
)

type GetDogsLinkGenerator interface {
	GenerateDogLink(dogId string) string
	GenerateShelterLink(shelterId string) string
	GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO
}

type GetDogsRepository interface {
	GetDogs(ctx context.Context, params dto.QueryParams) (dogdto.GetDogsQueryResponseDTO, error)
}

type GetDogsService struct {
	repo          GetDogsRepository
	linkGenerator GetDogsLinkGenerator
}

func NewGetDogsService(repo GetDogsRepository, linkGenerator GetDogsLinkGenerator) GetDogsService {
	return GetDogsService{
		repo:          repo,
		linkGenerator: linkGenerator,
	}
}

func (g GetDogsService) GetDogs(ctx context.Context, queryParams dto.QueryParams) (dogdto.DogsAndPaginationLinksDTO, error) {
	emptyDto := dogdto.DogsAndPaginationLinksDTO{}
	dogsResult, err := g.repo.GetDogs(ctx, queryParams)
	if err != nil {
		return emptyDto, err
	}

	var dogs []dogdto.DogDTO
	for _, dog := range dogsResult.Dogs {
		dogJson, err := json.Marshal(dog.ToJson())
		if err != nil {
			return emptyDto, err
		}
		var dogDto dogdto.DogDTO
		err = json.Unmarshal(dogJson, &dogDto)
		if err != nil {
			return emptyDto, err
		}
		selfLink := g.linkGenerator.GenerateDogLink(fmt.Sprintf("%d", dog.Id))
		shelterLink := g.linkGenerator.GenerateShelterLink(fmt.Sprintf("%d", dog.ShelterId))
		dogDto.Links = dogdto.DogLinksDTO{
			SelfLink:    selfLink,
			ShelterLink: shelterLink,
		}
		dogs = append(dogs, dogDto)
	}
	paginationLinks := g.linkGenerator.GeneratePaginationLinks(dogsResult.TotalAmountAvailable, queryParams, "/dogs")
	dogsAndPaginationLinks := dogdto.DogsAndPaginationLinksDTO{
		Dogs:            dogs,
		PaginationLinks: paginationLinks,
	}
	return dogsAndPaginationLinks, nil
}
