package middleware

import (
	"go-fiber-api/internal/util/token"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(c *fiber.Ctx) error {
	tokenStr := c.Cookies("token")

	if tokenStr == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized: missing auth token",
		})
	}

	claims, err := token.ParseToken(tokenStr)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized: invalid auth token",
		})
	}

	c.Locals("user_id", claims)

	return c.Next()
}
