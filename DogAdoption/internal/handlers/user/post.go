package userhandler

import (
	userdto "1dv027/aad/internal/dto/user"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type PostUserService interface {
	CreateNewUser(ctx context.Context, newUser userdto.NewUserDTO) (userdto.UserDTO, error)
}

type PostUserHandler struct {
	service PostUserService
}

func NewPostUserHandler(service PostUserService) PostUserHandler {
	return PostUserHandler{
		service: service,
	}
}

// Handle creates a new user.
// @Summary Create a new user
// @Description Adds a new user to the system with the provided user data in JSON format.
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body      userdto.NewUserDTO  true  "New User Data"  "The information for creating a new user"
// @Success 201  {object}  userdto.UserDTO  "Success, returns the newly created user information"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the request body is incomplete or contains invalid data"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if something goes wrong internally"
// @Router /users [post]
func (p PostUserHandler) Handle(c *fiber.Ctx) error {
	var newUserDto userdto.NewUserDTO
	err := c.BodyParser(&newUserDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request. check documentation for endpoint information.",
		})
	}
	registeredUser, err := p.service.CreateNewUser(c.Context(), newUserDto)
	if err != nil {
		var incompleteNewUser *customerrors.IncompleteNewUserError
		if errors.As(err, &incompleteNewUser) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "incomplete request body. check documentation for endpoint information.",
			})
		}
		var invalidNewUserData *customerrors.InvalidNewUserDataError
		if errors.As(err, &invalidNewUserData) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body. check documentation for endpoint information.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(registeredUser)
}
