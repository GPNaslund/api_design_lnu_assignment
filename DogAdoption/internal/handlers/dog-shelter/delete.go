package dogshelterhandler

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type DeleteDogShelterService interface {
	DeleteDogShelter(ctx context.Context, shelterId string, credentials dto.UserCredentials) error
}

type DeleteDogShelterHandler struct {
	service DeleteDogShelterService
}

func NewDeleteShelterHandler(service DeleteDogShelterService) DeleteDogShelterHandler {
	return DeleteDogShelterHandler{
		service: service,
	}
}

// Handle deletes a dog shelter by ID.
// @Summary Delete a dog shelter
// @Description Deletes a dog shelter specified by its ID if the requester has the necessary permissions.
// @Tags dogshelters
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "Shelter ID"  "The unique identifier of the dog shelter to delete"
// @Success 204  "Dog shelter deleted successfully"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request if the request was malformed"
// @Failure 404  {object}  dto.ErrorResponse "Not Found if the dog shelter with the specified ID does not exist"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized if the user does not have permission to delete the dog shelter"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error for any server errors"
// @Router /dogshelters/{id} [delete]
// @Security BearerAuth
func (d DeleteDogShelterHandler) Handle(c *fiber.Ctx) error {
	userCredentials := c.Locals("user").(dto.UserCredentials)
	err := d.service.DeleteDogShelter(c.Context(), c.Params("id"), userCredentials)
	if err != nil {
		var notFoundError *customerrors.DogShelterNotFoundError
		if errors.As(err, &notFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "resource not found",
			})
		}
		var unauthorizedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		var integerConversionErr *customerrors.IntegerConversionError
		if errors.As(err, &integerConversionErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "id param must be a number",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later.",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
