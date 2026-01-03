package auth

import (
	"go-fiber-api/internal/common/response"
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
		return response.Fail(c, fiber.StatusInternalServerError, "failed to parse request body")
	}

	if errorMessages, err := validation.ValidateStruct(req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "validation internal error")
	} else if len(errorMessages) > 0 {
		return response.FailWithData(c, fiber.StatusBadRequest, "validation failed", errorMessages)
	}

	err := h.authService.RegisterUser(&req)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to register user")
	}

	return response.SuccessWithStatus[any](c, fiber.StatusCreated, "user registered successfully", nil)
}

func (h *authHandler) LoginUser(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to parse request body")
	}

	if errorMessages, err := validation.ValidateStruct(req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "validation internal error")
	} else if len(errorMessages) > 0 {
		return response.FailWithData(c, fiber.StatusBadRequest, "validation failed", errorMessages)
	}

	user, err := h.authService.LoginUser(&req)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to login user")
	}

	_, err = token.GenerateToken(c, user.ID)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to generate token")
	}

	return response.Success(c, "login successful", user)
}

func (h *authHandler) GetUser(c *fiber.Ctx) error {
	v := c.Locals("user_id")
	if v == nil {
		return response.Fail(c, fiber.StatusUnauthorized, "unauthorized: missing user id")
	}

	claims, ok := v.(*token.CustomClaims)
	if !ok || claims == nil {
		return response.Fail(c, fiber.StatusUnauthorized, "unauthorized: invalid user claims")
	}

	user, err := h.authService.GetUser(claims.UserID)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to get user")
	}

	return response.Success(c, "get user success", user)

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

	return response.Success[any](c, "logout successful", nil)
}
