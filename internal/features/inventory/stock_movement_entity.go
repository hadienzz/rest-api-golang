package inventory

import (
	"time"

	"github.com/google/uuid"
)

type StockMovementType string

const (
	StockIn     StockMovementType = "IN"
	StockOut    StockMovementType = "OUT"
	StockAdjust StockMovementType = "ADJUST"
	StockSale   StockMovementType = "SALE"
)

type StockMovement struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`

	Type     StockMovementType `gorm:"type:varchar(10);not null"`
	Quantity int               `gorm:"type:int;not null"`

	ReferenceID   *uuid.UUID `gorm:"type:uuid;index"`
	ReferenceType string     `gorm:"type:varchar(50);index"` // TRANSACTION, RESTOCK, ADJUSTMENT

	CreatedAt time.Time
}
