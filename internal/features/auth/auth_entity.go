package auth

import (
	"database/sql"
	"go-fiber-api/internal/features/merchant"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`

	Merchants []merchant.Merchant `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
