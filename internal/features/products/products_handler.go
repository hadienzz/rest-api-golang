package products

import (
	"context"
	"fmt"
	"go-fiber-api/internal/common/response"
	"go-fiber-api/internal/util/token"
	"go-fiber-api/internal/util/upload"
	"go-fiber-api/internal/util/validation"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler interface {
	CreateProduct(c *fiber.Ctx) error
	GetMerchantProducts(c *fiber.Ctx) error
	BulkDeleteMerchantProducts(c *fiber.Ctx) error
	GetMerchantProductsDashboard(c *fiber.Ctx) error
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
	// Ambil user dari token dan form-data
	form, err := c.MultipartForm()
	ctx, cancel := context.WithTimeout(c.Context(), 45*time.Second)
	defer cancel()
	log.Print(form)
	if err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid form data")
	}

	user_id := c.Locals("user_id").(*token.CustomClaims).UserID
	merchantIDParam := c.Params("merchant_id")

	var request CreateProductRequest
	if err := c.BodyParser(&request); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Validasi field-field dasar
	if validationErrors, err := validation.ValidateStruct(request); err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "validation error")
	} else if len(validationErrors) > 0 {
		return response.FailWithData(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	// Validasi merchant_id
	merchantID, err := uuid.Parse(merchantIDParam)
	if err != nil || merchantID == uuid.Nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid merchant id")
	}

	// Ambil merchant
	merchant, err := ph.merchantService.GetMerchantById(merchantID)
	if err != nil {
		return response.Fail(c, fiber.StatusNotFound, "merchant not found")
	}

	// Authorization: pastikan merchant milik user
	if merchant.UserID != user_id {
		return response.Fail(c, fiber.StatusForbidden, "you are not allowed to access this merchant")
	}

	// Handle upload product photo
	productPhotoFiles := form.File["product_photo_url"]
	if len(productPhotoFiles) == 0 {
		return response.Fail(c, fiber.StatusBadRequest, "product photo is required")
	}

	uploadResult, err := upload.UploadToSupabaseStorage(ctx, productPhotoFiles[0], "products")
	if err != nil {
		return response.FailWithData(c, fiber.StatusInternalServerError, "failed to upload product photo", err.Error())
	}

	// Lengkapi request untuk service layer
	request.MerchantID = merchantID
	request.ProductPhotoUrl = uploadResult.PublicURL

	createdProduct, err := ph.productService.CreateProduct(&request)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to create product")
	}

	return response.SuccessWithStatus(c, fiber.StatusCreated, "product created", createdProduct)
}

func (ph *productHandler) GetMerchantProducts(c *fiber.Ctx) error {
	params := c.Params("id")

	merchantID, err := uuid.Parse(params)
	if err != nil || merchantID == uuid.Nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid merchant id")
	}

	products, err := ph.productService.GetMerchantProducts(merchantID)

	if err != nil {
		return response.FailWithData(c, fiber.StatusInternalServerError, "failed to get products", err.Error())
	}

	return response.Success(c, "products retrieved", products)
}

func (ph *productHandler) BulkDeleteMerchantProducts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(*token.CustomClaims).UserID

	var req BulkDeleteProductRequest
	fmt.Println(req.ProductIDs)
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	if len(req.ProductIDs) == 0 {
		return response.Fail(c, fiber.StatusBadRequest, "product_ids cannot be empty")
	}

	merchant, err := ph.merchantService.GetMyMerchants(userID)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to get merchant")
	}

	if merchant.UserID != userID {
		return response.Fail(c, fiber.StatusForbidden, "you are not allowed to access this merchant")
	}

	if err := ph.productService.DeleteMerchantProduct(
		req.ProductIDs,
		merchant.ID,
	); err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessNoData(c, "products deleted")
}

func (ph *productHandler) GetMerchantProductsDashboard(c *fiber.Ctx) error {
	params := c.Params("merchant_id")

	merchantID, err := uuid.Parse(params)

	if err != nil || merchantID == uuid.Nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid merchant id")
	}

	products, err := ph.productService.GetMerchantProductsDashboard(merchantID)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to get products")
	}

	return response.Success(c, "product dashboard received", products)

}
