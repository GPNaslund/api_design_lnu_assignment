package service

import (
	"1dv027/aad/internal/dto"
	"fmt"
	"net/url"
	"strings"
)

type HateoasLinkGenerator struct {
	basePath string
}

func NewHateoasLinkGenerator(basePath string) HateoasLinkGenerator {
	return HateoasLinkGenerator{
		basePath: basePath,
	}
}

func (d HateoasLinkGenerator) GenerateEntryPointLinks() dto.EntryPointLinksDTO {
	return dto.EntryPointLinksDTO{
		OpenApi:           fmt.Sprintf("%s/swagger", d.basePath),
		DogsUrl:           fmt.Sprintf("%s/dogs", d.basePath),
		DogSheltersUrl:    fmt.Sprintf("%s/dogshelters", d.basePath),
		AuthenticationUrl: fmt.Sprintf("%s/auth/login", d.basePath),
		UsersUrl:          fmt.Sprintf("%s/users", d.basePath),
	}
}

func (d HateoasLinkGenerator) GenerateDogLink(dogId string) string {
	return fmt.Sprintf("%s/dogs/%s", d.basePath, dogId)
}

func (d HateoasLinkGenerator) GenerateShelterLink(shelterId string) string {
	return fmt.Sprintf("%s/dogshelters/%s", d.basePath, shelterId)
}

func (d HateoasLinkGenerator) GenerateDogsFromDogShelterLink(shelterId string) string {
	return fmt.Sprintf("%s/dogs?shelter-id=%s", d.basePath, shelterId)
}

func (d HateoasLinkGenerator) GeneratePaginationLinks(totalItems int, queryParams dto.QueryParams, apiPath string) dto.PaginationLinksDTO {
	pageSize := *queryParams.Pagination.Limit
	currentPage := *queryParams.Pagination.Page

	totalPages := (totalItems + pageSize - 1) / pageSize
	links := dto.PaginationLinksDTO{}

	filters := d.getFiltersMapFromDto(queryParams)
	var filterParams []string
	for key, value := range filters {
		filterParams = append(filterParams, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
	}
	filterQueryString := strings.Join(filterParams, "&")

	base := fmt.Sprintf("%s%s?", d.basePath, apiPath)
	if len(filterQueryString) > 0 {
		base += filterQueryString + "&"
	}

	links.Self = fmt.Sprintf("%spage=%d&size=%d", base, currentPage, pageSize)

	links.First = fmt.Sprintf("%spage=1&size=%d", base, pageSize)
	links.Last = fmt.Sprintf("%spage=%d&size=%d", base, totalPages, pageSize)

	if currentPage > 1 {
		links.Prev = fmt.Sprintf("%spage=%d&size=%d", base, currentPage-1, pageSize)
	}

	if currentPage < totalPages {
		links.Next = fmt.Sprintf("%spage=%d&size=%d", base, currentPage+1, pageSize)
	}

	return links
}

func (d HateoasLinkGenerator) getFiltersMapFromDto(queryParams dto.QueryParams) map[string]string {
	var appliedFilters = make(map[string]string)
	dogsFilters := queryParams.DogsFilter
	if dogsFilters != nil {
		if dogsFilters.Gender != nil {
			appliedFilters["gender"] = *dogsFilters.Gender
		}
		if dogsFilters.Breed != nil {
			appliedFilters["breed"] = *dogsFilters.Breed
		}
		if dogsFilters.IsAdopted != nil {
			appliedFilters["is-adopted"] = *dogsFilters.IsAdopted
		}
		if dogsFilters.IsNeutered != nil {
			appliedFilters["is-neutered"] = *dogsFilters.IsNeutered
		}
		if dogsFilters.ShelterId != nil {
			appliedFilters["shelter-id"] = fmt.Sprintf("%d", *dogsFilters.ShelterId)
		}
	}
	dogShelterFilters := queryParams.DogShelterFilter
	if dogShelterFilters != nil {
		if dogShelterFilters.City != nil {
			appliedFilters["city"] = *dogShelterFilters.City
		}
		if dogShelterFilters.Country != nil {
			appliedFilters["country"] = *dogShelterFilters.Country
		}
		if dogShelterFilters.Name != nil {
			appliedFilters["name"] = *dogShelterFilters.Name
		}
	}

	return appliedFilters
}
