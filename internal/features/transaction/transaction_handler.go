package transaction

import (
	"log"
	"net/http"

	"go-fiber-api/internal/common/response"
	"go-fiber-api/internal/util/token"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	claimsValue := c.Locals("user_id")

	claims, ok := claimsValue.(*token.CustomClaims)
	if !ok || claims == nil {
		return response.Fail(c, http.StatusUnauthorized, "unauthorized")
	}

	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, http.StatusBadRequest, "invalid request body")
	}

	log.Println("Creating transaction merchant: ", req.MerchantID, "idempotencykey:", req.IdempotencyKey)

	result, err := h.service.CreateTransaction(claims.UserID, &req)
	if err != nil {
		// Midtrans library kadang mengembalikan error interface dengan
		// dynamic type *midtrans.Error namun pointer-nya nil, sehingga
		// pemanggilan err.Error() akan menyebabkan panic (nil receiver).
		if midErr, ok := err.(*midtrans.Error); ok && midErr == nil {
			return response.Fail(c, http.StatusBadRequest, "payment gateway error")
		}

		return response.Fail(c, http.StatusBadRequest, err.Error())
	}

	return response.Success(c, "transaction cd", result)
}

func (h *TransactionHandler) HandleMidtransWebhook(c *fiber.Ctx) error {

	var notification MidtransNotificationRequest
	if err := c.BodyParser(&notification); err != nil {
		log.Println("invalid payload:", err)
		return c.SendStatus(http.StatusOK)
	}

	if err := h.service.HandleMidtransWebhook(&notification); err != nil {
		log.Println("webhook error:", err)
		// tetap 200
		return c.SendStatus(http.StatusOK)
	}

	return c.SendStatus(http.StatusOK)
}

func (h *TransactionHandler) GetTransactionDetail(c *fiber.Ctx) error {
	orderID := c.Params("orderId")
	if orderID == "" {
		return response.Fail(c, http.StatusBadRequest, "order_id is required")
	}

	result, err := h.service.GetTransactionDetail(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Fail(c, http.StatusNotFound, "transaction not found")
		}
		return response.Fail(c, http.StatusBadRequest, err.Error())
	}

	return response.Success(c, "transaction detail", result)
}

func (h *TransactionHandler) GetTransactionsByUserID(c *fiber.Ctx) error {
	user_id := c.Locals("user_id").(*token.CustomClaims).UserID

	result, err := h.service.GetTransactionsByUserID(user_id)
	if err != nil {
		return response.Fail(c, http.StatusBadRequest, err.Error())

	}

	return response.Success(c, "user transactions", result)
}
