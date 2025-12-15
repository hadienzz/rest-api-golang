package auth

import (
	"go-fiber-api/internal/util/token"
	"go-fiber-api/internal/util/validation"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	authService AuthService
}

type Handler interface {
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
}

func NewHandler(service AuthService) *authHandler {
	return &authHandler{
		authService: service,
	}
}

func (h *authHandler) RegisterUser(c *fiber.Ctx) error {
	var req RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to parse request body",
		})
	}

	if errorMessages, err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation internal error",
		})
	} else if len(errorMessages) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Failed",
			"errors":  errorMessages,
		})
	}

	err := h.authService.RegisterUser(&req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
	})
}

func (h *authHandler) LoginUser(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to parse request body",
		})
	}

	if errorMessages, err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation internal error",
		})
	} else if len(errorMessages) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Failed",
			"errors":  errorMessages,
		})
	}

	user, err := h.authService.LoginUser(&req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to login user",
		})
	}

	_, err = token.GenerateToken(c, user.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    user,
	})
}

func (h *authHandler) GetUser(c *fiber.Ctx) error {
	v := c.Locals("user_id")
	if v == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized: missing user id",
		})
	}

	claims, ok := v.(*token.CustomClaims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized: invalid user claims",
		})
	}

	user, err := h.authService.GetUser(claims.UserID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get user success",
		"data":    user,
	})

}

func (h *authHandler) LogoutUser(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "", // kosongkan
		Path:     "/",
		Expires:  time.Unix(0, 0), // waktu lampau
		MaxAge:   -1,              // beberapa browser honor ini
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
		// Domain: "example.com", // hanya set kalau kamu set waktu create cookie
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
