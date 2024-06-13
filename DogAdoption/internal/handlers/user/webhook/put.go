package userwebhookhandler

import (
	"1dv027/aad/internal/dto"
	userwebhookdto "1dv027/aad/internal/dto/user/webhook"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type PutUserWebhookService interface {
	UpdateUserWebhook(ctx context.Context, idParam string, user dto.UserCredentials, data userwebhookdto.UpdateUserWebhookDTO) (userwebhookdto.UserWebhookDTO, error)
}

type PutUserWebhookHandler struct {
	service PutUserWebhookService
}

func NewPutUserWebhookHandler(service PutUserWebhookService) PutUserWebhookHandler {
	return PutUserWebhookHandler{
		service: service,
	}
}

// Handle updates a specific user webhook by ID.
// @Summary Update a user webhook
// @Description Updates information for a specific user webhook identified by its unique ID for the authenticated user based on the provided data in JSON format. Secret must be minimum 12 characters.
// @Tags users/{id}/webhook
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "User ID"  "The unique identifier of the user of which webhook to update"
// @Param   data  body      userwebhookdto.UpdateUserWebhookDTO  true  "Update Webhook Data"  "The updated information for the user webhook"
// @Success 200  {object}  userwebhookdto.UserWebhookDTO  "Success, returns the updated user webhook information"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the request data is incomplete or has invalid values"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user credentials do not match or are invalid"
// @Failure 404  {object}  dto.ErrorResponse "Not Found, if no matching webhook resource is found for the given ID"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /users/{id}/webhook [put]
// @Security BearerAuth
func (p PutUserWebhookHandler) Handle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userCredentials := c.Locals("user").(dto.UserCredentials)
	var updateWebhookDto userwebhookdto.UpdateUserWebhookDTO
	err := c.BodyParser(&updateWebhookDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later.",
		})
	}

	webhookDto, err := p.service.UpdateUserWebhook(c.Context(), idParam, userCredentials, updateWebhookDto)
	if err != nil {
		var unauthorizedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		var incompleteWebhookDataError *customerrors.IncompleteWebhookDataError
		if errors.As(err, &incompleteWebhookDataError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "bad data in request body. visit documentation for more endpoint information.",
			})
		}
		var invalidWebhookDataError *customerrors.InvalidWebhookDataError
		if errors.As(err, &invalidWebhookDataError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid values in request body. visit documentation for more endpoint information.",
			})
		}
		var webhookNotFoundError *customerrors.WebhookNotFoundError
		if errors.As(err, &webhookNotFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "resource not found.",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(webhookDto)
}
