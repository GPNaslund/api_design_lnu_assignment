package doghandler

import (
	"1dv027/aad/internal/dto"
	dogdto "1dv027/aad/internal/dto/dog"
	customerrors "1dv027/aad/internal/errors"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type PostDogService interface {
	CreateDog(ctx context.Context, newDogDto dogdto.NewDogDTO, credentials dto.UserCredentials) (dogdto.DogDTO, error)
}

type PostDogHandler struct {
	service PostDogService
}

func NewPostDogHandler(service PostDogService) PostDogHandler {
	return PostDogHandler{
		service: service,
	}
}

// Handle posts a new dog to the system.
// @Summary Add a new dog
// @Description Adds a new dog to the system with the provided dog data in JSON format. shelter_id field is for admins only.
// @Tags dogs
// @Accept  json
// @Produce  json
// @Param   dog   body      dogdto.NewDogDTO  true  "Dog Data"  "The dog information to be created"
// @Success 201  {object}  dogdto.DogDTO  "Success, returns the newly created dog information"
// @Failure 400  {object}  dto.ErrorResponse "Bad Request, if the JSON body cannot be parsed or mandatory fields are missing"
// @Failure 401  {object}  dto.ErrorResponse "Unauthorized, if the user does not have permission to add a dog"
// @Failure 500  {object}  dto.ErrorResponse "Internal Server Error, if an error occurs while processing the request"
// @Router /dogs [post]
// @Security BearerAuth
func (p PostDogHandler) Handle(c *fiber.Ctx) error {
	var newDogDto dogdto.NewDogDTO
	err := c.BodyParser(&newDogDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "bad request body. visit the documentation for endpoint information.",
		})
	}

	userCredentials := c.Locals("user").(dto.UserCredentials)

	createdDog, err := p.service.CreateDog(c.Context(), newDogDto, userCredentials)
	if err != nil {
		var incompleteDogDataErr *customerrors.IncompleteDogDataError
		if errors.As(err, &incompleteDogDataErr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "request body is incomplete. visit documentation for more information.",
			})
		}
		var unauthorizedErr *customerrors.UnauthorizedError
		if errors.As(err, &unauthorizedErr) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "something went wrong internally. try again later!",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createdDog)
}
