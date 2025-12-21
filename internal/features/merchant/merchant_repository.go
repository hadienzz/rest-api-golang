package merchant

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	CreateMerchant(merchant *Merchant) (*Merchant, error)
	GetMerchantById(id uuid.UUID) (*Merchant, error)
	GetAllMerchant() ([]Merchant, error)
	GetMyMerchant(userID uuid.UUID) (*Merchant, error) // nanti ganti jadi dashboard
	GetMyMerchantsSummary(userID uuid.UUID) ([]MerchantSummary, error)
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

func (mr *merchantRepository) CreateMerchant(merchant *Merchant) (*Merchant, error) {
	if err := mr.db.Create(merchant).Error; err != nil {
		return nil, err
	}
	return merchant, nil
}

func (mr *merchantRepository) GetMerchantById(id uuid.UUID) (*Merchant, error) {
	var merchant Merchant

	err := mr.db.
		Where("id = ?", id).
		First(&merchant).
		Error

	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (mr *merchantRepository) GetAllMerchant() ([]Merchant, error) {
	var merchants []Merchant
	result := mr.db.Find(&merchants)
	return merchants, result.Error
}

func (mr *merchantRepository) GetMyMerchant(
	userID uuid.UUID,
) (*Merchant, error) {

	var merchant Merchant

	err := mr.db.
		Where("user_id = ?", userID).
		Take(&merchant).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (mr *merchantRepository) GetMyMerchantsSummary(
	userID uuid.UUID,
) ([]MerchantSummary, error) {

	var merchant []MerchantSummary

	err := mr.db.
		Table("merchants").
		Select("id, user_id, name, description, profile_photo_url").
		Where("user_id = ?", userID).
		Find(&merchant).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return merchant, nil
}
