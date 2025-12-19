package products

import (
	"go-fiber-api/internal/util/token"
	"go-fiber-api/internal/util/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler interface {
	CreateProduct(c *fiber.Ctx) error
	GetMerchantProducts(c *fiber.Ctx) error
}

type productHandler struct {
	productService  ProductService
	merchantService MerchantServiceContract
}

func NewProductHandler(productService ProductService, merchantService MerchantServiceContract) ProductHandler {
	return &productHandler{
		productService:  productService,
		merchantService: merchantService,
	}
}

func (ph *productHandler) CreateProduct(c *fiber.Ctx) error {
	// Ambil user dari token

	user_id := c.Locals("user_id").(*token.CustomClaims).UserID
	MerchantID := c.Params("merchant_id")

	var request CreateProductRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	// Validasi request (API Contract)
	if validationErrors, err := validation.ValidateStruct(request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "validation error",
		})
	} else if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation failed",
			"errors":  validationErrors,
		})
	}

	// Validasi merchant_id
	merchantID, err := uuid.Parse(MerchantID)
	if err != nil || merchantID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid merchant id",
		})
	}

	// Ambil merchant
	merchant, err := ph.merchantService.GetMerchantById(merchantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "merchant not found",
		})
	}

	// Authorization: pastikan merchant milik user
	if merchant.UserID != user_id {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "you are not allowed to access this merchant",
		})
	}

	request.MerchantID = merchantID
	createdProduct, err := ph.productService.CreateProduct(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "product created",
		"data":    createdProduct,
	})
}

func (ph *productHandler) GetMerchantProducts(c *fiber.Ctx) error {
	params := c.Params("id")

	products, err := ph.productService.GetMerchantProducts(params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get products",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "products retrieved",
		"data":    products,
	})
}
