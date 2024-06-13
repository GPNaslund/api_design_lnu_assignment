package userwebhookhandler

import (
	"1dv027/aad/internal/dto"
	userwebhookdto "1dv027/aad/internal/dto/user/webhook"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type PostUserWebhookService interface {
	CreateNewWebhook(ctx context.Context, idParam string, user dto.UserCredentials, webhookData userwebhookdto.NewUserWebhookDTO) (userwebhookdto.UserWebhookDTO, error)
}

type PostUserWebhookHandler struct {
	service PostUserWebhookService
}

func NewPostUserWebhookHandler(service PostUserWebhookService) PostUserWebhookHandler {
	return PostUserWebhookHandler{
		service: service,
	}
}

// Handle creates a new webhook for the user.
// @Summary Create a new user webhook
// @Description Adds a new webhook for the authenticated user based on the provided webhook data in JSON format. Secret must be minimum 12 characters.
// @Tags users/{id}/webhook
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param   id   path      integer  true  "User ID"  "The unique identifier of the user for whom the webhook is being created"
// @Param   webhook  body      userwebhookdto.NewUserWebhookDTO  true  "Webhook Data"  "The information for the new user webhook"
// @Success 201  {object}  userwebhookdto.UserWebhookDTO  "Success, returns the newly created user webhook information"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the JSON body cannot be parsed, mandatory fields are missing, or the webhook data is incomplete"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user credentials do not match or are invalid"
// @Failure 404  {object}  dto.ErrorResponse "Not Found, if the specified user does not exist"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /users/{id}/webhook [post]
func (p PostUserWebhookHandler) Handle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userCredentials := c.Locals("user").(dto.UserCredentials)
	var newWebhookDto userwebhookdto.NewUserWebhookDTO
	err := c.BodyParser(&newWebhookDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request. visit documentation for more endpoint information.",
		})
	}

	webhookDto, err := p.service.CreateNewWebhook(c.Context(), idParam, userCredentials, newWebhookDto)
	if err != nil {
		var incompleteWebhookDataError *customerrors.IncompleteWebhookDataError
		if errors.As(err, &incompleteWebhookDataError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "bad request data. visit documentation for more endpoint information.",
			})
		}
		var invalidWebhookDataError *customerrors.InvalidWebhookDataError
		if errors.As(err, &invalidWebhookDataError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "bad request data. visit documentation for more endpoint information.",
			})
		}
		var userNotFoundError *customerrors.UserNotFoundError
		if errors.As(err, &userNotFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "resource not found.",
			})
		}
		var unauthorziedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorziedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(webhookDto)
}
