package userwebhookhandler

import (
	"1dv027/aad/internal/dto"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type DeleteUserWebhookService interface {
	DeleteWebhook(ctx context.Context, idParam string, user dto.UserCredentials) error
}

type DeleteUserWebhookHandler struct {
	service DeleteUserWebhookService
}

func NewDeleteUserWebhookHandler(service DeleteUserWebhookService) DeleteUserWebhookHandler {
	return DeleteUserWebhookHandler{
		service: service,
	}
}

// Handle deletes a webhook for the authenticated user.
// @Summary Delete a webhook
// @Description Deletes a webhook for the authenticated user based on the provided webhook ID.
// @Tags users/{id}/webhook
// @Accept  json
// @Produce  json
// @Param   id   path      integer  true  "Webhook ID"  "The unique identifier of the webhook to delete"
// @Success 204  "Webhook deleted successfully"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request - No webhook found or ID must be a number"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized - Invalid user credentials"
// @Failure 404  {object}  dto.ErrorResponse "Not Found - No webhook found"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error - Something went wrong internally, try again later"
// @Router /users/{id}/webhook [delete]
// @Security BearerAuth
func (d DeleteUserWebhookHandler) Handle(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userCredentials := c.Locals("user").(dto.UserCredentials)
	err := d.service.DeleteWebhook(c.Context(), idParam, userCredentials)
	if err != nil {
		var webhookNotFound *customerrors.WebhookNotFoundError
		if errors.As(err, &webhookNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "no webhook found",
			})
		}
		var unauthorizedError *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedError) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		var integerConversionError *customerrors.IntegerConversionError
		if errors.As(err, &integerConversionError) {
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
