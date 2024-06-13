package doghandler

import (
	dogdto "1dv027/aad/internal/dto/dog"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type GetDogByIdService interface {
	GetDogById(ctx context.Context, dogId string) (dogdto.DogDTO, error)
}

type GetDogByIdHandler struct {
	service GetDogByIdService
}

func NewGetDogByIdHandler(service GetDogByIdService) GetDogByIdHandler {
	return GetDogByIdHandler{
		service: service,
	}
}

// Handle retrieves a dog by its ID.
// @Summary Get a dog by ID
// @Description Retrieves detailed information about a dog specified by its ID.
// @Tags dogs
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "Dog ID"  "The unique identifier of the dog to retrieve"
// @Success 200  {object}  dogdto.DogDTO  "Success, returns detailed information about the dog"
// @Failure 404  {object}  dto.ErrorResponse "Not Found, if no dog matches the provided ID"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /dogs/{id} [get]
func (g GetDogByIdHandler) Handle(c *fiber.Ctx) error {
	dogDataResponse, err := g.service.GetDogById(c.Context(), c.Params("id"))
	if err != nil {
		var notFoundErr *customerrors.DogNotFoundError
		if errors.As(err, &notFoundErr) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "dog not found",
			})
		}
		var integerConversionErr *customerrors.IntegerConversionError
		if errors.As(err, &integerConversionErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid id parameter",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusOK).JSON(dogDataResponse)
}
