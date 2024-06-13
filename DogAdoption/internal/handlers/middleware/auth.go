package middleware

import (
	"1dv027/aad/internal/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type JwtService interface {
	ValidateToken(jwt string) (dto.UserCredentials, error)
}

type AuthMiddleware struct {
	jwtService JwtService
}

func NewAuthMiddleware(jwtService JwtService) AuthMiddleware {
	return AuthMiddleware{
		jwtService: jwtService,
	}
}

func (a AuthMiddleware) AuthenticateRequest(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "authorization header is required",
		})
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}

	token := bearerToken[1]

	userData, err := a.jwtService.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	c.Locals("user", userData)
	return c.Next()
}
