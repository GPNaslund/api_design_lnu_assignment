package dogshelterhandler

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	"context"

	"github.com/gofiber/fiber/v2"
)

type GetDogShelterService interface {
	GetShelters(ctx context.Context, queryParams dto.QueryParams) (dogshelterdto.DogSheltersAndPaginationLinksDTO, error)
}

type GetDogShelterHandler struct {
	service GetDogShelterService
}

func NewGetDogShelterHandler(service GetDogShelterService) GetDogShelterHandler {
	return GetDogShelterHandler{
		service: service,
	}
}

// Handle retrieves dog shelters based on query parameters.
// @Summary Get dog shelters
// @Description Retrieves a list of dog shelters based on provided query parameters like location and capacity.
// @Tags dogshelters
// @Accept  json
// @Produce  json
// @Param   name   query     string  false  "Filter by name"
// @Param   country   query     string  false  "Filter by country"
// @Param   city   query     string     false  "Filter by city"
// @Success 200  {object}  dogshelterdto.DogSheltersAndPaginationLinksDTO  "Success, returns a list of dog shelters"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request if the query parameters are invalid"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error if an error occurs while processing the request"
// @Router /dogshelters [get]
func (g GetDogShelterHandler) Handle(c *fiber.Ctx) error {
	queryParamsDto := c.Locals("queryParams").(dto.QueryParams)

	getDogSheltersResult, err := g.service.GetShelters(c.Context(), queryParamsDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusOK).JSON(getDogSheltersResult)
}
