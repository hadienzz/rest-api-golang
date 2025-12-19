package merchant

import (
	"go-fiber-api/internal/features/products"

	"github.com/google/uuid"
)

type MerchantServiceAdapter struct {
	service MerchantService
}

func NewMerchantServiceAdapter(service MerchantService) products.MerchantServiceContract {
	return &MerchantServiceAdapter{
		service: service,
	}
}

func (a *MerchantServiceAdapter) GetMerchantById(
	id uuid.UUID,
) (*products.MerchantInfo, error) {

	merchant, err := a.service.GetMerchantById(id)

	if err != nil {
		return nil, err
	}

	return &products.MerchantInfo{
		ID:              merchant.ID,
		UserID:          merchant.UserID,
		Name:            merchant.Name,
		Description:     merchant.Description,
		ProfilePhotoUrl: merchant.ProfilePhotoUrl,
		BannerImageUrl:  merchant.BannerImageUrl,
	}, nil
}

func (a *MerchantServiceAdapter) GetMyMerchants(userID uuid.UUID) (*products.MerchantInfo, error) {

	merchants, err := a.service.GetMyMerchant(userID)

	if err != nil {
		return nil, err
	}

	return &products.MerchantInfo{
		ID:              merchants.ID,
		UserID:          merchants.UserID,
		Name:            merchants.Name,
		Description:     merchants.Description,
		ProfilePhotoUrl: merchants.ProfilePhotoUrl,
	}, nil
}
