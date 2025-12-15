package products

import (
	"go-fiber-api/internal/util/token"
	"go-fiber-api/internal/util/validation"

	"github.com/gofiber/fiber/v2"
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
	// Ambil user_id dari token
	userIDFromToken := c.Locals("user_id").(*token.CustomClaims).UserID

	// Cek dulu apakah user ini punya merchant
	merchants, err := ph.merchantService.GetMyMerchant(userIDFromToken)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get merchant",
			"error":   err.Error(),
		})
	}

	if len(merchants) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user does not have merchant",
		})
	}

	// Ambil merchant pertama (atau bisa disesuaikan kalau multi-merchant)
	merchantID := merchants[0].ID.String()

	// Parse body request produk
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	// Set merchant_id dari merchant yang ditemukan, user tidak boleh kirim sendiri
	req.MerchantID = merchantID

	// Validasi payload
	if errorMessages, err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation internal error",
		})
	} else if len(errorMessages) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation failed",
			"errors":  errorMessages,
		})
	}

	// Mapping ke entity Product
	product := &Product{
		MerchantID:  req.MerchantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	createdProduct, err := ph.productService.CreateProduct(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create product",
			"error":   err.Error(),
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
