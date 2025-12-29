package products

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MerchantID uuid.UUID `gorm:"type:uuid;not null;index"`

	Name            string          `gorm:"type:varchar(100);not null"`
	Description     string          `gorm:"type:text"`
	Price           decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	Quantity        int             `gorm:"not null"`
	ProductPhotoUrl string          `gorm:"type:text;not null"`

	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
