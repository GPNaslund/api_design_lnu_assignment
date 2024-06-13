package dogshelterhandler

import (
	"1dv027/aad/internal/dto"
	dogshelterdto "1dv027/aad/internal/dto/dog-shelter"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type PostDogShelterService interface {
	CreateDogShelter(ctx context.Context, newDogShelter dogshelterdto.NewDogShelterDTO, credentials dto.UserCredentials) (dogshelterdto.DogShelterDTO, error)
}

type PostDogShelterHandler struct {
	service PostDogShelterService
}

func NewPostShelterHandler(service PostDogShelterService) PostDogShelterHandler {
	return PostDogShelterHandler{
		service: service,
	}
}

// Handle creates a new dog shelter.
// @Summary Add a new dog shelter
// @Description Adds a new dog shelter to the system with the provided shelter data in JSON format.
// @Tags dogshelters
// @Accept  json
// @Produce  json
// @Param   shelter  body      dogshelterdto.NewDogShelterDTO  true  "Dog Shelter Data"  "The information for the new dog shelter"
// @Success 201  {object}  dogshelterdto.DogShelterDTO  "Success, returns the newly created dog shelter information"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the JSON body cannot be parsed, mandatory fields are missing, or the dog shelter data is incomplete"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user does not have permission to add a dog shelter"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /dogshelters [post]
// @Security BearerAuth
func (p PostDogShelterHandler) Handle(c *fiber.Ctx) error {
	var newDogShelterDto dogshelterdto.NewDogShelterDTO
	err := c.BodyParser(&newDogShelterDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request body. visit documentation for endpoint information.",
		})
	}
	userCredentials := c.Locals("user").(dto.UserCredentials)

	newDogShelter, err := p.service.CreateDogShelter(c.Context(), newDogShelterDto, userCredentials)
	if err != nil {
		var incompleteDogShelterDataErr *customerrors.IncompleteDogShelterDataError
		if errors.As(err, &incompleteDogShelterDataErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "request body incomplete. visit documentation for endpoint information.",
			})
		}
		var unauthorizedErr *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedErr) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		var invalidDogShelterDataErr *customerrors.InvalidNewDogShelterDataError
		if errors.As(err, &invalidDogShelterDataErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "request body is not valid. visit documentation for endpoint information.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(newDogShelter)
}
