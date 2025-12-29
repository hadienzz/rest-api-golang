package merchant

import (
	"context"
	"go-fiber-api/internal/common/response"
	"go-fiber-api/internal/util/token"
	"go-fiber-api/internal/util/upload"
	"go-fiber-api/internal/util/validation"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MerchantHandler interface {
	AddMerchant(c *fiber.Ctx) error
	GetMerchantById(c *fiber.Ctx) error
	GetAllMerchant(c *fiber.Ctx) error
	GetMyMerchantDashboard(c *fiber.Ctx) error
	GetMyMerchantsSummary(c *fiber.Ctx) error
	GetMerchantDisplay(c *fiber.Ctx) error
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
	form, err := c.MultipartForm()
	if err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "failed to parse multipart form")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 45*time.Second)
	defer cancel()

	// =========================
	// 1. PROFILE PHOTO (single)
	// =========================
	profileFiles := form.File["profile_photo_url"]
	if len(profileFiles) == 0 {
		return response.Fail(c, fiber.StatusBadRequest, "profile_photo is required")
	}

	profileResult, err := upload.UploadToSupabaseStorage(
		ctx,
		profileFiles[0],
		"profiles",
	)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to upload profile photo")
	}

	// =========================
	// 2. BANNER IMAGE (single)
	// =========================
	bannerFiles := form.File["banner_image_url"]
	var bannerURL string

	if len(bannerFiles) > 0 {
		bannerResult, err := upload.UploadToSupabaseStorage(
			ctx,
			bannerFiles[0],
			"banners",
		)
		if err != nil {
			return response.Fail(c, fiber.StatusInternalServerError, "failed to upload banner image")
		}
		bannerURL = bannerResult.PublicURL
	}

	// =========================
	// 3. GALLERY PHOTOS (multiple)
	// =========================
	galleryFiles := form.File["gallery_photos"]
	if len(galleryFiles) > 6 {
		return response.Fail(c, fiber.StatusBadRequest, "maximum 6 gallery photos allowed")
	}

	galleryURLs := make([]string, 0, len(galleryFiles))

	for _, fileHeader := range galleryFiles {
		result, err := upload.UploadToSupabaseStorage(
			ctx,
			fileHeader,
			"galleries",
		)
		if err != nil {
			return response.Fail(c, fiber.StatusInternalServerError, "failed to upload gallery photo")
		}
		galleryURLs = append(galleryURLs, result.PublicURL)
	}

	// =========================
	// 4. BUILD REQUEST DTO
	// =========================
	userIDFromToken := c.Locals("user_id").(*token.CustomClaims).UserID

	latitudeStr := c.FormValue("latitude")
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid latitude value")
	}

	longitudeStr := c.FormValue("longitude")
	longitude, err := strconv.ParseFloat(longitudeStr, 64)

	if err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid longitude value")
	}

	req := MerchantDTO{
		UserID:          userIDFromToken,
		Name:            c.FormValue("name"),
		Description:     c.FormValue("description"),
		Type:            c.FormValue("type"),
		Location:        c.FormValue("location"),
		ProfilePhotoUrl: profileResult.PublicURL,
		BannerImageUrl:  bannerURL,
		GalleryPhotoUrl: galleryURLs,
		GoogleMapUrl:    c.FormValue("google_maps_url"),
		IFrameMapUrl:    c.FormValue("iframe_maps_url"),
		Latitude:        latitude,
		Longitude:       longitude,
	}

	log.Println("MerchantDTO:", req)

	// =========================
	// 5. VALIDATION
	// =========================
	if errorMessages, err := validation.ValidateStruct(req); err != nil {
		log.Println("validation internal error:", err)
		return response.Fail(c, fiber.StatusBadRequest, "validation internal error")
	} else if len(errorMessages) > 0 {
		log.Println("validation failed:", errorMessages)
		return response.FailWithData(c, fiber.StatusBadRequest, "validation failed", errorMessages)
	}

	// =========================
	// 6. SERVICE CALL
	// =========================
	if err := h.merchantService.AddMerchant(&req); err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to create merchant")
	}

	return response.SuccessWithStatus[any](c, fiber.StatusCreated, "merchant created", nil)
}

func (h *merchantHandler) GetMerchantById(c *fiber.Ctx) error {
	merchantId := c.Params("id")

	merchantUUID, err := uuid.Parse(merchantId)

	if err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "invalid merchant id format")
	}

	merchant, err := h.merchantService.GetMerchantById(merchantUUID)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "merchant not found")
	}

	return response.Success(c, "merchant retrieved", merchant)
}

func (h *merchantHandler) GetAllMerchant(c *fiber.Ctx) error {
	merchants, err := h.merchantService.GetAllMerchant()

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to retrieve merchants")
	}

	return response.Success(c, "merchants retrieved", merchants)
}

func (h *merchantHandler) GetMyMerchantDashboard(c *fiber.Ctx) error {
	claims := c.Locals("user_id").(*token.CustomClaims)
	userID := claims.UserID

	merchant, err := h.merchantService.GetMyMerchant(userID)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to retrieve merchant")
	}

	if merchant == nil {
		return response.Success[any](c, "merchant not found", nil)
	}

	return response.Success(c, "merchant retrieved", merchant)
}

func (h *merchantHandler) GetMyMerchantsSummary(c *fiber.Ctx) error {
	user_id := c.Locals("user_id").(*token.CustomClaims).UserID

	merchant, err := h.merchantService.GetMyMerchantsSummary(user_id)

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to retrieve merchants")
	}

	if len(merchant) == 0 {
		return response.Success(c, "this account does not have merchant", []MerchantSummary{})
	}

	return response.Success(c, "merchants retrieved", merchant)
}

func (h *merchantHandler) GetMerchantDisplay(c *fiber.Ctx) error {
	merchants, err := h.merchantService.GetMerchantDisplay()

	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "failed to get merchant display")
	}

	return response.Success(c, "merchants retrieved", merchants)
}
