package merchant

import (
	"fmt"
	"go-fiber-api/internal/features/products"

	"github.com/google/uuid"
)

type MerchantService interface {
	AddMerchant(req *CreateMerchantRequest) error
	GetMerchantById(id uuid.UUID) (*MerchantDTO, error)
	GetAllMerchant() ([]MerchantDTO, error)
	GetMyMerchant(userID uuid.UUID) ([]products.MerchantInfo, error)
}

type merchantService struct {
	merchantRepository MerchantRepository
}

func NewMerchantService(merchantRepo MerchantRepository) MerchantService {
	return &merchantService{
		merchantRepository: merchantRepo,
	}
}

func (ms *merchantService) AddMerchant(req *CreateMerchantRequest) error {
	merchant := &Merchant{
		UserID:       req.UserID,
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		Location:     req.Location,
		ProfilePhoto: req.ProfilePhoto,
		// BannerImg:    req.BannerImg.Filename,
	}

	_, err := ms.merchantRepository.CreateMerchant(merchant)
	return err
}

func (ms *merchantService) GetMerchantById(id uuid.UUID) (*MerchantDTO, error) {
	merchant, err := ms.merchantRepository.GetMerchantById(id)

	result := &MerchantDTO{
		ID:           merchant.ID,
		UserID:       merchant.UserID,
		Name:         merchant.Name,
		Description:  merchant.Description,
		Type:         merchant.Type,
		Location:     merchant.Location,
		ProfilePhoto: merchant.ProfilePhoto,
		CreatedAt:    merchant.CreatedAt,
		UpdatedAt:    merchant.UpdatedAt,
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
			ID:           m.ID,
			UserID:       m.UserID,
			Name:         m.Name,
			Description:  m.Description,
			Type:         m.Type,
			Location:     m.Location,
			ProfilePhoto: m.ProfilePhoto,
		}
	}

	return result, nil
}
func (ms *merchantService) GetMyMerchant(userID uuid.UUID) ([]products.MerchantInfo, error) {
	merchants, err := ms.merchantRepository.GetMyMerchant(userID)
	if err != nil {
		return nil, err
	}

	result := make([]products.MerchantInfo, len(merchants))
	for i, m := range merchants {
		// mapping produk entity -> DTO
		result[i] = products.MerchantInfo{
			ID:           m.ID,
			UserID:       m.UserID,
			Name:         m.Name,
			Description:  m.Description,
			Type:         m.Type,
			Location:     m.Location,
			ProfilePhoto: m.ProfilePhoto,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		}
	}

	return result, nil
}
