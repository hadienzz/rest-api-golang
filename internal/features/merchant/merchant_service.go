package merchant

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type MerchantService interface {
	AddMerchant(req *MerchantDTO) error
	GetMerchantById(id uuid.UUID) (*MerchantDTO, error)
	GetAllMerchant() ([]MerchantDTO, error)
	GetMyMerchant(userID uuid.UUID) (*MerchantDTO, error)
	GetMyMerchantsSummary(userID uuid.UUID) ([]MerchantSummary, error)
	GetMerchantDisplay() ([]MerchantSummary, error)
}

type merchantService struct {
	merchantRepository MerchantRepository
}

func NewMerchantService(merchantRepo MerchantRepository) MerchantService {
	return &merchantService{
		merchantRepository: merchantRepo,
	}
}

func (ms *merchantService) AddMerchant(req *MerchantDTO) error {
	merchant := &Merchant{
		UserID:          req.UserID,
		Name:            req.Name,
		Description:     req.Description,
		Type:            req.Type,
		Location:        req.Location,
		ProfilePhotoUrl: req.ProfilePhotoUrl,
		BannerImageUrl:  req.BannerImageUrl,
		GalleryPhotoUrl: req.GalleryPhotoUrl,
		GoogleMapUrl:    req.GoogleMapUrl,
		IFrameMapUrl:    req.IFrameMapUrl,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
	}

	_, err := ms.merchantRepository.CreateMerchant(merchant)

	if err != nil {
		return err
	}

	return nil
}

func (ms *merchantService) GetMerchantById(id uuid.UUID) (*MerchantDTO, error) {
	merchant, err := ms.merchantRepository.GetMerchantById(id)

	result := &MerchantDTO{
		ID:              merchant.ID,
		UserID:          merchant.UserID,
		Name:            merchant.Name,
		Description:     merchant.Description,
		Type:            merchant.Type,
		Location:        merchant.Location,
		ProfilePhotoUrl: merchant.ProfilePhotoUrl,
		CreatedAt:       merchant.CreatedAt,
		UpdatedAt:       merchant.UpdatedAt,
		BannerImageUrl:  merchant.BannerImageUrl,
		GalleryPhotoUrl: merchant.GalleryPhotoUrl,
		Restricted:      merchant.Restricted,
		Verified:        merchant.Verified,
		TotalFollowers:  merchant.TotalFollowers,
		GoogleMapUrl:    merchant.GoogleMapUrl,
		IFrameMapUrl:    merchant.IFrameMapUrl,
		Latitude:        merchant.Latitude,
		Longitude:       merchant.Longitude,
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get merchant by id: %w", err.Error())
	}

	return result, err
}

func (ms *merchantService) GetAllMerchant() ([]MerchantDTO, error) {
	merchants, err := ms.merchantRepository.GetAllMerchant()

	if err != nil {
		return nil, err
	}

	result := make([]MerchantDTO, len(merchants))
	for i, m := range merchants {
		result[i] = MerchantDTO{
			ID:              m.ID,
			UserID:          m.UserID,
			Name:            m.Name,
			Description:     m.Description,
			Type:            m.Type,
			Location:        m.Location,
			ProfilePhotoUrl: m.ProfilePhotoUrl,
		}
	}

	return result, nil
}

func (ms *merchantService) GetMyMerchant(
	userID uuid.UUID,
) (*MerchantDTO, error) {

	merchant, err := ms.merchantRepository.GetMyMerchant(userID)
	if err != nil {
		return nil, err
	}

	if merchant == nil {
		return nil, nil
	}

	return &MerchantDTO{
		ID:              merchant.ID,
		UserID:          merchant.UserID,
		Name:            merchant.Name,
		Description:     merchant.Description,
		Type:            merchant.Type,
		Location:        merchant.Location,
		ProfilePhotoUrl: merchant.ProfilePhotoUrl,
		BannerImageUrl:  merchant.BannerImageUrl,
		GalleryPhotoUrl: merchant.GalleryPhotoUrl,
		CreatedAt:       merchant.CreatedAt,
		UpdatedAt:       merchant.UpdatedAt,
	}, nil
}

func (ms *merchantService) GetMyMerchantsSummary(userID uuid.UUID) ([]MerchantSummary, error) {
	merchants, err := ms.merchantRepository.GetMyMerchantsSummary(userID)

	if err != nil {
		return nil, err
	}

	result := make([]MerchantSummary, len(merchants))

	for i, m := range merchants {
		result[i] = MerchantSummary{
			ID:              m.ID,
			UserID:          m.UserID,
			Name:            m.Name,
			Description:     m.Description,
			ProfilePhotoUrl: m.ProfilePhotoUrl,
		}
	}

	return result, nil
}

func (ms *merchantService) GetMerchantDisplay() ([]MerchantSummary, error) {
	merchants, err := ms.merchantRepository.GetMerchantDisplay()
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	result := make([]MerchantSummary, len(merchants))
	for i, m := range merchants {
		result[i] = MerchantSummary{
			ID:              m.ID,
			Name:            m.Name,
			Description:     m.Description,
			ProfilePhotoUrl: m.ProfilePhotoUrl,
		}
	}

	return result, nil
}
