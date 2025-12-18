package products

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MerchantID  uuid.UUID       `gorm:"column:merchant_id;not null;index"`
	Name        string          `gorm:"column:name;type:varchar(100);not null"`
	Description string          `gorm:"column:description;type:text"`
	Price       decimal.Decimal `gorm:"column:price;type:decimal(10,2);not null"`
	Quantity    int             `gorm:"column:quantity;type:int;not null"`

	CreatedAt sql.NullTime `gorm:"column:created_at"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at"`
}
