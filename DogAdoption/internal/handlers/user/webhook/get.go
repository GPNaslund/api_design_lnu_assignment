package userwebhookhandler

import (
	"1dv027/aad/internal/dto"
	userwebhookdto "1dv027/aad/internal/dto/user/webhook"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type GetUserWebhookService interface {
	GetUserWebhook(ctx context.Context, idParam string, userCredentials dto.UserCredentials) (userwebhookdto.UserWebhookDTO, error)
}

type GetUserWebhookHandler struct {
	service GetUserWebhookService
}

func NewGetUserWebhookHandler(service GetUserWebhookService) GetUserWebhookHandler {
	return GetUserWebhookHandler{
		service: service,
	}
}

// Handle retrieves a specific user webhook by ID.
// @Summary Get a user webhook
// @Description Retrieves detailed information about a specific user webhook identified by its unique ID for the authenticated user.
// @Tags users/{id}/webhook
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "User ID"  "The unique identifier of the user to retrieve its webhook"
// @Success 200  {object}  userwebhookdto.UserWebhookDTO  "Success, returns detailed information about the user webhook"
// @Failure 404  {object}  dto.ErrorResponse "Not Found, if no resource matches the provided ID or the webhook does not belong to the user"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user credentials are invalid or do not grant access to the requested resource"
// @Router /users/{id}/webhook [get]
// @Security BearerAuth
func (g GetUserWebhookHandler) Handle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userCredentials := c.Locals("user").(dto.UserCredentials)

	webhookDto, err := g.service.GetUserWebhook(c.Context(), idParam, userCredentials)
	if err != nil {
		var noWebhookFoundError *customerrors.WebhookNotFoundError
		if errors.As(err, &noWebhookFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "no resource found",
			})
		}
		var unauthorizedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(webhookDto)
}
