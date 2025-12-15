package merchant

import (
	"database/sql"

	"github.com/google/uuid"
)

type CreateMerchantRequest struct {
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Type         string  `json:"type" validate:"required"`
	Location     float32 `json:"location" validate:"required"`
	ProfilePhoto string  `json:"profile_photo" validate:"required,url"` // ✅ simpan URL, bukan file
	UserID       string  `json:"-"`                                     // diisi dari token
}

type MerchantDTO struct {
	ID           uuid.UUID `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         string    `json:"type"`
	Location     float32   `json:"location"`
	ProfilePhoto string    `json:"profile_photo"`

	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
