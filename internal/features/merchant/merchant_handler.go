package merchant

import (
	"context"
	"go-fiber-api/internal/util/token"
	"go-fiber-api/internal/util/upload"
	"go-fiber-api/internal/util/validation"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MerchantHandler interface {
	AddMerchant(c *fiber.Ctx) error
	GetMerchantById(c *fiber.Ctx) error
	GetAllMerchant(c *fiber.Ctx) error
	GetMyMerchant(c *fiber.Ctx) error
}

type merchantHandler struct {
	merchantService MerchantService
}

func NewMerchantHandler(service MerchantService) MerchantHandler {
	return &merchantHandler{
		merchantService: service,
	}
}

func (h *merchantHandler) AddMerchant(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")
	if contentType == "" || !strings.Contains(contentType, "multipart/form-data") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "request must be multipart/form-data",
		})
	}

	var req CreateMerchantRequest

	// 1. Ambil file dari form
	fileHeader, err := c.FormFile("profile_photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "photo is required",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 45*time.Second)
	defer cancel()

	// 2. Upload ke Supabase pakai helper
	uploadResult, err := upload.UploadToSupabaseStorage(
		&ctx,
		fileHeader,
		"profiles",
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "upload to supabase failed",
			"error":   err.Error(),
		})
	}

	// 3. Ambil field text
	nameValue := c.FormValue("name")
	descValue := c.FormValue("description")
	typeValue := c.FormValue("type")
	locationValue := c.FormValue("location")

	parsedLocation, err := strconv.ParseFloat(locationValue, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "location must be a valid number",
			"error":   err.Error(),
		})
	}

	userIDFromToken := c.Locals("user_id").(*token.CustomClaims).UserID

	// 4. Isi struct request
	req.Name = nameValue
	req.Description = descValue
	req.Type = typeValue
	req.Location = float32(parsedLocation)
	req.UserID = userIDFromToken.String()
	req.ProfilePhoto = uploadResult.PublicURL
	// 5. Validasi
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

	if err := h.merchantService.AddMerchant(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create merchant",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "merchant created",
	})
}

func (h *merchantHandler) GetMerchantById(c *fiber.Ctx) error {
	merchantId := c.Params("id")

	merchantUUID, err := uuid.Parse(merchantId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid merchant id format",
		})
	}

	merchant, err := h.merchantService.GetMerchantById(merchantUUID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "merchant not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "merchant retrieved",
		"data":    merchant,
	})
}

func (h *merchantHandler) GetAllMerchant(c *fiber.Ctx) error {
	merchants, err := h.merchantService.GetAllMerchant()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve merchants",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "merchants retrieved",
		"data":    merchants,
	})
}

func (h *merchantHandler) GetMyMerchant(c *fiber.Ctx) error {
	userIDFromToken := c.Locals("user_id").(*token.CustomClaims).UserID

	merchants, err := h.merchantService.GetMyMerchant(userIDFromToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrieve merchants",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "merchants retrieved",
		"data":    merchants,
	})
}
