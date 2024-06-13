package usermehandler

import (
	"1dv027/aad/internal/dto"
	userdto "1dv027/aad/internal/dto/user"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type GetUserMeService interface {
	GetAuthenticatedUser(ctx context.Context, userCredentials dto.UserCredentials) (userdto.UserDTO, error)
}

type GetUserMeHandler struct {
	service GetUserMeService
}

func NewGetUserMeHandler(service GetUserMeService) GetUserMeHandler {
	return GetUserMeHandler{
		service: service,
	}
}

// Handle retrieves the authenticated user's information.
// @Summary Get authenticated user information
// @Description Retrieves detailed information about the authenticated user based on the provided credentials.
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200  {object}  userdto.UserDTO  "Success, returns detailed information about the authenticated user"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user credentials do not match or are invalid"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /users/me [get]
func (g GetUserMeHandler) Handle(c *fiber.Ctx) error {
	userCredentials := c.Locals("user").(dto.UserCredentials)

	userDto, err := g.service.GetAuthenticatedUser(c.Context(), userCredentials)
	if err != nil {
		var unauthorizedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": userDto,
	})
}
