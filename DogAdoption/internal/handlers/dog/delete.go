package doghandler

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type DeleteDogService interface {
	DeleteDog(ctx context.Context, dogId string, credentials dto.UserCredentials) error
}

type DeleteDogHandler struct {
	service DeleteDogService
}

func NewDeleteDogHandler(service DeleteDogService) DeleteDogHandler {
	return DeleteDogHandler{
		service: service,
	}
}

// Handle deletes a dog by ID.
// @Summary Delete a dog
// @Description Deletes a dog specified by its ID if the requester has the necessary permissions.
// @Tags dogs
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "Dog ID"
// @Success 204  "Dog deleted successfully"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request if the request was malformed"
// @Failure 404  {object}  dto.ErrorResponse "Not Found if the dog with the specified ID does not exist"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized if the user does not have permission to delete the dog"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error for any server errors"
// @Router /dogs/{id} [delete]
// @Security BearerAuth
func (d DeleteDogHandler) Handle(c *fiber.Ctx) error {
	userCredentials := c.Locals("user").(dto.UserCredentials)
	err := d.service.DeleteDog(c.Context(), c.Params("id"), userCredentials)
	if err != nil {
		var notFoundError *customerrors.DogNotFoundError
		if errors.As(err, &notFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "dog not found found",
			})
		}
		var unauthorizedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		var IntegerConversionError *customerrors.IntegerConversionError
		if errors.As(err, &IntegerConversionError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "id parameter must be a number",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. Try again later.",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
