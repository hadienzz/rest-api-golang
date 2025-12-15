package merchant

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	CreateMerchant(merchant *Merchant) (*Merchant, error)
	GetMerchantById(id uuid.UUID) (*Merchant, error)
	GetAllMerchant() ([]Merchant, error)
	GetMyMerchant(userID uuid.UUID) ([]Merchant, error)
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
	result := mr.db.Create(merchant)
	if result.Error != nil {
		log.Println("ERROR: ", result.Error)
	}
	return merchant, result.Error
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

func (mr *merchantRepository) GetMyMerchant(userID uuid.UUID) ([]Merchant, error) {
	var merchants []Merchant
	err := mr.db.Where("user_id = ?", userID).Find(&merchants).Error

	return merchants, err
}
