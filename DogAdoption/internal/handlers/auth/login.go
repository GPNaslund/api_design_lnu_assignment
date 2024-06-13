package authhandler

import (
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Token string `json:"token"`
}

type AuthService interface {
	ValidateUsernameAndPassword(ctx context.Context, username, password string) (string, error)
	GetAllowedFields() map[string]any
}

type RequestBodyValidator interface {
	ValidateRequestBody(allowedFields map[string]any, requestBody []byte) error
}

type LoginHandler struct {
	authService   AuthService
	bodyValidator RequestBodyValidator
}

func NewLoginHandler(authService AuthService, bodyValidator RequestBodyValidator) LoginHandler {
	return LoginHandler{
		authService:   authService,
		bodyValidator: bodyValidator,
	}
}

// HandleLogin logs in a user and returns a JWT token.
// @Summary User login
// @Description Authenticates a user by username and password, and returns a JWT token if successful.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   payload  body      Payload  true  "Login Credentials"
// @Success 200  {object}  Response "Returns JWT token"
// @Failure 400  {object}  dto.ErrorResponse "Bad request when the JSON body cannot be parsed or wrong payload type"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, when the username or password is incorrect"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, something went wrong with the server"
// @Router /auth/login [post]
func (l LoginHandler) Handle(c *fiber.Ctx) error {
	allowedFields := l.authService.GetAllowedFields()
	err := l.bodyValidator.ValidateRequestBody(allowedFields, c.Body())
	if err != nil {
		// For invalid request bodies
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var payload Payload
	if err = c.BodyParser(&payload); err != nil {
		// If wrong content type
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "wrong payload type. read documentation for more information.",
		})
	}

	jwt, err := l.authService.ValidateUsernameAndPassword(c.Context(), payload.Username, payload.Password)
	if err != nil {
		// Specific error handling
		var wrongCredentialsErr *customerrors.WrongCredentialsError
		var unauthorizedErr *customerrors.UnauthorizedError
		if errors.As(err, &wrongCredentialsErr) || errors.As(err, &unauthorizedErr) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid username or password",
			})
		}
		// For all other errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong with the server. try again later",
		})
	}

	// Successful response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": jwt,
	})

}
