package middleware

import (
	"1dv027/aad/internal/dto"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type QueryParamsValidator struct {
}

func NewQueryParamsValidator() QueryParamsValidator {
	return QueryParamsValidator{}
}

func (q QueryParamsValidator) ValidateQueryParams(c *fiber.Ctx) error {
	queryParams := c.Queries()
	paginationParams, err := q.validatePaginationParams(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	dogFilterParams, err := q.validateDogsParams(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	dogShelterParams, err := q.validateDogShelterParams(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(queryParams) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query params supplied. Please check documentation for valid query params.",
		})
	}

	queryParamsDto := dto.QueryParams{
		Pagination:       &paginationParams,
		DogsFilter:       &dogFilterParams,
		DogShelterFilter: &dogShelterParams,
	}

	c.Locals("queryParams", queryParamsDto)

	return c.Next()

}

func (q QueryParamsValidator) validatePaginationParams(query map[string]string) (dto.PaginationParams, error) {
	paginationParams := dto.PaginationParams{}
	page := query["page"]
	limit := query["limit"]

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return paginationParams, fmt.Errorf("invalid page variable")
		}
		paginationParams.Page = &pageInt
		delete(query, "page")
	} else {
		defaultPage := 1
		paginationParams.Page = &defaultPage
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return paginationParams, fmt.Errorf("invalid limit variable")
		}
		paginationParams.Limit = &limitInt
		delete(query, "limit")
	} else {
		defaultLimit := 10
		paginationParams.Limit = &defaultLimit
	}

	return paginationParams, nil
}

func (q QueryParamsValidator) validateDogsParams(query map[string]string) (dto.DogsFilterParams, error) {
	dogFilterParams := dto.DogsFilterParams{}

	breed := query["breed"]
	gender := query["gender"]
	isNeutered := query["is-neutered"]
	isAdopted := query["is-adopted"]
	shelterId := query["shelter-id"]

	if breed != "" {
		if len(breed) > 64 {
			return dogFilterParams, fmt.Errorf("breed name is to long")
		}
		dogFilterParams.Breed = &breed
		delete(query, "breed")
	}

	if gender != "" {
		if gender != "male" && gender != "female" {
			return dogFilterParams, fmt.Errorf("invalid gender value. only male and female are accepted")
		}
		dogFilterParams.Gender = &gender
		delete(query, "gender")
	}

	if isNeutered != "" {
		if isNeutered != "true" && isNeutered != "false" {
			return dogFilterParams, fmt.Errorf("invalid is_neutered value. only true and false are accepted")
		}
		dogFilterParams.IsNeutered = &isNeutered
		delete(query, "is-neutered")
	}

	if isAdopted != "" {
		if isAdopted != "true" && isAdopted != "false" {
			return dogFilterParams, fmt.Errorf("invalid is_adopted value. only true and false are accepted")
		}
		dogFilterParams.IsAdopted = &isAdopted
		delete(query, "is-adopted")
	}

	if shelterId != "" {
		shelterIdInt, err := strconv.Atoi(shelterId)
		if err != nil {
			return dogFilterParams, fmt.Errorf("shelter-id must be a number")
		}
		dogFilterParams.ShelterId = &shelterIdInt
	}
	return dogFilterParams, nil
}

func (q QueryParamsValidator) validateDogShelterParams(query map[string]string) (dto.DogShelterFilterParams, error) {
	dogShelterFilterParams := dto.DogShelterFilterParams{}

	country := query["country"]
	city := query["city"]
	name := query["name"]

	if country != "" {
		if len(country) > 65 {
			return dogShelterFilterParams, fmt.Errorf("value for country is too long")
		}
		dogShelterFilterParams.Country = &country
		delete(query, "country")
	}

	if city != "" {
		if len(city) > 65 {
			return dogShelterFilterParams, fmt.Errorf("value for city is too long")
		}
		dogShelterFilterParams.City = &city
		delete(query, "city")
	}

	if name != "" {
		if len(name) > 65 {
			return dogShelterFilterParams, fmt.Errorf("value for name is too long")
		}
		dogShelterFilterParams.Name = &name
		delete(query, "name")
	}

	return dogShelterFilterParams, nil
}
