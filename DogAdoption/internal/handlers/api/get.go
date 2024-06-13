package apihandler

import (
	"1dv027/aad/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type ApiService interface {
	GetEntryPointLinks() dto.EntryPointLinksDTO
}

type ApiHandler struct {
	service ApiService
}

func NewApiHandler(service ApiService) ApiHandler {
	return ApiHandler{
		service: service,
	}
}

// GetEntryPointLinks returns the entry point links for the API.
// @Summary Get entry point links
// @Description Returns a collection of links that represent the entry point of the API.
// @Tags entrypoint
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.EntryPointLinksDTO
// @Router / [get]
func (a ApiHandler) Handle(c *fiber.Ctx) error {
	entryPointLinks := a.service.GetEntryPointLinks()
	return c.Status(fiber.StatusOK).JSON(entryPointLinks)
}
