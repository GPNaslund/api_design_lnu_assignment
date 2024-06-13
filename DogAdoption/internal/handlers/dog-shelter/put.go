package dogshelterhandler

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type PutDogShelterService interface {
	UpdateDogShelter(ctx context.Context,
		dogId string, credentials dto.UserCredentials,
		updateData dogshelterdto.UpdateDogShelterDTO) (dogshelterdto.DogShelterDTO, error)
}

type PutDogShelterHandler struct {
	service PutDogShelterService
}

func NewPutDogShelterHandler(service PutDogShelterService) PutDogShelterHandler {
	return PutDogShelterHandler{
		service: service,
	}
}

// Handle updates an existing dog shelter by ID.
// @Summary Update a dog shelter
// @Description Updates the information for an existing dog shelter specified by its ID with the provided shelter data in JSON format.
// @Tags dogshelters
// @Accept  json
// @Produce  json
// @Param   id   path      string  true  "Shelter ID"  "The unique identifier of the dog shelter to update"
// @Param   shelter  body      dogshelterdto.UpdateDogShelterDTO  true  "Dog Shelter Update Data"  "The dog shelter information updates to apply"
// @Success 200  {object}  dogshelterdto.DogShelterDTO  "Success, returns the updated dog shelter information"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the JSON body cannot be parsed, mandatory fields are missing, or the dog shelter data is incomplete"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user does not have permission to update the dog shelter"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /dogshelters/{id} [put]
// @Security BearerAuth
func (p PutDogShelterHandler) Handle(c *fiber.Ctx) error {
	userCredentials := c.Locals("user").(dto.UserCredentials)
	var updateDogShelterDto dogshelterdto.UpdateDogShelterDTO
	err := c.BodyParser(&updateDogShelterDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request. check documentation for endpoint information",
		})
	}

	updatedDogShelter, err := p.service.UpdateDogShelter(c.Context(), c.Params("id"), userCredentials, updateDogShelterDto)
	if err != nil {
		var unauthorizedErr *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedErr) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var incompleteDogDataErr *customerrors.IncompleteDogShelterDataError
		if errors.As(err, &incompleteDogDataErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "bad request body. check documentation for endpoint information",
			})
		}
		var integerConversionError *customerrors.IntegerConversionError
		if errors.As(err, &integerConversionError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "id param must be a number",
			})
		}
		var dogShelterNotFoundError *customerrors.DogShelterNotFoundError
		if errors.As(err, &dogShelterNotFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "resource not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusOK).JSON(updatedDogShelter)
}
