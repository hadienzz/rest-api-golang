package follow

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	MerchantID uuid.UUID `gorm:"type:uuid;not null;index"`

	CreatedAt time.Time
}
