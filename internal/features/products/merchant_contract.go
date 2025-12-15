package products

import (
	"database/sql"

	"github.com/google/uuid"
)

// DTO untuk product yang dikirim ke client

// DTO untuk merchant yang dipakai di response
type MerchantInfo struct {
	ID           uuid.UUID    `json:"id"`
	UserID       string       `json:"user_id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Type         string       `json:"type"`
	Location     float32      `json:"location"`
	ProfilePhoto string       `json:"profile_photo"`
	CreatedAt    sql.NullTime `json:"created_at"`
	UpdatedAt    sql.NullTime `json:"updated_at"`
}

type MerchantServiceContract interface {
	GetMyMerchant(userID uuid.UUID) ([]MerchantInfo, error)
}
