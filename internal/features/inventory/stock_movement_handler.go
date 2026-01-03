package inventory

import (
	"go-fiber-api/internal/common/response"

	"github.com/gofiber/fiber/v2"
)

type StockMovementHandler interface {
	AddStockIn(c *fiber.Ctx) error
	// AddStockOut(c *fiber.Ctx) error
}

type stockMovementHandler struct {
	service StockMovementService
}

func NewStockMovementHandler(service StockMovementService) StockMovementHandler {
	return &stockMovementHandler{
		service: service,
	}
}

func (h *stockMovementHandler) AddStockIn(c *fiber.Ctx) error {
	var request StockMovementDTO
	if err := c.BodyParser(&request); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	err := h.service.AddStockIn(request.ProductID, request.Quantity)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessNoData(c, "stock in added successfully")
}
