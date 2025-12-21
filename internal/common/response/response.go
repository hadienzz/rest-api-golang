package response

import "github.com/gofiber/fiber/v2"

type ApiResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func Success[T any](c *fiber.Ctx, message string, data T) error {
	return c.Status(fiber.StatusOK).JSON(ApiResponse[T]{
		Message: message,
		Data:    data,
	})
}

func Fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ApiResponse[any]{
		Message: message,
	})
}

func SuccessWithStatus[T any](c *fiber.Ctx, status int, message string, data T) error {
	return c.Status(status).JSON(ApiResponse[T]{
		Message: message,
		Data:    data,
	})
}

func FailWithData[T any](c *fiber.Ctx, status int, message string, data T) error {
	return c.Status(status).JSON(ApiResponse[T]{
		Message: message,
		Data:    data,
	})
}

func SuccessNoData(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}
