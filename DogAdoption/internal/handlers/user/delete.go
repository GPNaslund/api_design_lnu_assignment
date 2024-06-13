package userhandler

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type DeleteUserService interface {
	DeleteUser(ctx context.Context, idParam string, user dto.UserCredentials) error
}

type DeleteUserHandler struct {
	service DeleteUserService
}

func NewDeleteUserHandler(service DeleteUserService) DeleteUserHandler {
	return DeleteUserHandler{
		service: service,
	}
}

// Handle deletes a specific user by ID.
// @Summary Delete a user
// @Description Deletes a user identified by its unique ID, provided the requester has the necessary permissions.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "User ID"  "The unique identifier of the user to delete"
// @Success 204  "User deleted successfully"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the ID parameter format is incorrect"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the requester does not have permission to delete the user"
// @Failure 404  {object}  dto.ErrorResponse "Not Found, if no user matches the provided ID"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /users/{id} [delete]
// @Security BearerAuth
func (d DeleteUserHandler) Handle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userCredentials := c.Locals("user").(dto.UserCredentials)

	err := d.service.DeleteUser(c.Context(), idParam, userCredentials)
	if err != nil {
		var userNotFound *customerrors.UserNotFoundError
		if errors.As(err, &userNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		var integerConversionError *customerrors.IntegerConversionError
		if errors.As(err, &integerConversionError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "id parameter must be a number",
			})
		}

		var unauthorizedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized to perform action",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later.",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
