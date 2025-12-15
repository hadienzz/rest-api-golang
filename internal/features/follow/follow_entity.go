package follow

import (
	"database/sql"

	"github.com/google/uuid"
)

type Follow struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"column:user_id;not null;index"`
	MerchantID uuid.UUID `gorm:"column:merchant_id;type:uuid;not null;index"`

	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
