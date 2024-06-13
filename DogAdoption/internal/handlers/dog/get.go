package doghandler

import (
	"1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	"context"

	"github.com/gofiber/fiber/v2"
)

type GetDogsService interface {
	GetDogs(ctx context.Context, queryParams dto.QueryParams) (dogdto.DogsAndPaginationLinksDTO, error)
}

type GetDogsHandler struct {
	service GetDogsService
}

func NewGetDogsHandler(service GetDogsService) GetDogsHandler {
	return GetDogsHandler{
		service: service,
	}
}

// Handle retrieves dogs based on query parameters.
// @Summary Get dogs
// @Description Retrieves a list of dogs based on provided query parameters like breed, size, and age.
// @Tags dogs
// @Accept  json
// @Produce  json
// @Param   breed   query     string  false  "Filter by dog breed"
// @Param   gender    query     string  false  "Filter by dog gender"
// @Param   is_neutered     query     boolean     false  "Filter by if dog is neutered"
// @Param   is_adopted     query     boolean     false  "Filter by if dog is adopted"
// @Param shelter_id	query	integer		false	"Filter dogs that are from a specific dog shelter"
// @Success 200  {object}  dogdto.DogsAndPaginationLinksDTO  "Success, returns a list of dogs along with pagination details"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request if the query parameters are invalid"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error if an error occurs while processing the request"
// @Router /dogs [get]
func (g GetDogsHandler) Handle(c *fiber.Ctx) error {
	queryParamsDto := c.Locals("queryParams").(dto.QueryParams)

	dogsDataResponse, err := g.service.GetDogs(c.Context(), queryParamsDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request",
		})
	}
	return c.Status(fiber.StatusOK).JSON(dogsDataResponse)
}
