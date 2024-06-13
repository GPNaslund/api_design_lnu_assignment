package dogshelterhandler

import (
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type GetDogShelterByIdService interface {
	GetDogShelterById(ctx context.Context, dogId string) (dogshelterdto.DogShelterDTO, error)
}

type GetDogShelterByIdHandler struct {
	service GetDogShelterByIdService
}

func NewGetDogShelterByIdHandler(service GetDogShelterByIdService) GetDogShelterByIdHandler {
	return GetDogShelterByIdHandler{
		service: service,
	}
}

// Handle retrieves a specific dog shelter by ID.
// @Summary Get a dog shelter by ID
// @Description Retrieves detailed information about a specific dog shelter identified by its unique ID.
// @Tags dogshelters
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "Shelter ID"  "The unique identifier of the dog shelter to retrieve"
// @Success 200  {object}  dogshelterdto.DogShelterDTO  "Success, returns detailed information about the dog shelter"
// @Failure 404  {object}  dto.ErrorResponse "Not Found, if no dog shelter matches the provided ID"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /dogshelters/{id} [get]
func (g GetDogShelterByIdHandler) Handle(c *fiber.Ctx) error {
	dogShelterData, err := g.service.GetDogShelterById(c.Context(), c.Params("id"))
	if err != nil {
		var notFoundErr *customerrors.DogShelterNotFoundError
		if errors.As(err, &notFoundErr) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "resource not found",
			})
		}
		var integerConversionErr *customerrors.IntegerConversionError
		if errors.As(err, &integerConversionErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "id param must be a number",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusOK).JSON(dogShelterData)
}
