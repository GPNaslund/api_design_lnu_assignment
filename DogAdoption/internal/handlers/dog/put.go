package doghandler

import (
	"1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type PutDogService interface {
	UpdateDog(ctx context.Context, dogId string, credentials dto.UserCredentials, updateData dogdto.UpdateDogDTO) (dogdto.DogDTO, error)
}

type PutDogHandler struct {
	service PutDogService
}

func NewPutDogHandler(service PutDogService) PutDogHandler {
	return PutDogHandler{
		service: service,
	}
}

// Handle updates an existing dog by ID.
// @Summary Update dog information
// @Description Updates the information for an existing dog specified by its ID with the provided dog data in JSON format.
// @Tags dogs
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "Dog ID"  "The unique identifier of the dog to update"
// @Param   dog   body      dogdto.UpdateDogDTO  true  "Dog Update Data"  "The dog information updates to apply"
// @Success 200  {object}  dogdto.DogDTO  "Success, returns the updated dog information"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the JSON body cannot be parsed, mandatory fields are missing, or the dog data is incomplete"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user does not have permission to update the dog"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /dogs/{id} [put]
// @Security BearerAuth
func (p PutDogHandler) Handle(c *fiber.Ctx) error {
	userCredentials := c.Locals("user").(dto.UserCredentials)
	var updateDogDto dogdto.UpdateDogDTO
	err := c.BodyParser(&updateDogDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request body. read the documentation for endpoint information.",
		})
	}

	updatedDog, err := p.service.UpdateDog(c.Context(), c.Params("id"), userCredentials, updateDogDto)
	if err != nil {
		var unauthorizedErr *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedErr) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		var incompleteDogDataErr *customerrors.IncompleteDogDataError
		if errors.As(err, &incompleteDogDataErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body data. visit the documentation for endpoint information.",
			})
		}
		var dogNotFoundError *customerrors.DogNotFoundError
		if errors.As(err, &dogNotFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "dog not found",
			})
		}
		var integerConversionError *customerrors.IntegerConversionError
		if errors.As(err, &integerConversionError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "id parameter must be a number",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusOK).JSON(updatedDog)
}
