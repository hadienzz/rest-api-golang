package transaction

import (
	"go-fiber-api/internal/features/products"

	"github.com/google/uuid"
)

type TransactionItem struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index"`

	Quantity int   `gorm:"type:int;not null"`
	Price    int64 `gorm:"type:bigint;not null"`
	Subtotal int64 `gorm:"type:bigint;not null"`

	// Relations
	Transaction Transaction      `gorm:"foreignKey:TransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Product     products.Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
