package merchant

import (
	"database/sql"
	"go-fiber-api/internal/features/products"

	"github.com/google/uuid"
)

type Merchant struct {
	ID           uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       string             `gorm:"column:user_id;not null;index"`
	Name         string             `gorm:"column:name;type:varchar(100);not null"`
	Description  string             `gorm:"column:description;type:text"`
	Type         string             `gorm:"column:type;type:varchar(50)"`
	Location     float32            `gorm:"column:location"`
	ProfilePhoto string             `gorm:"column:profile_photo;type:text"`
	Product      []products.Product `gorm:"foreignKey:MerchantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// BannerImg  string  `db:"banner_img"`   // kalau nanti pakai banner
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
