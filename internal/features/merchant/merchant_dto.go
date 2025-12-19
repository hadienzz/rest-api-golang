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
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Type            string    `json:"type"`
	Location        string    `json:"location"`
	ProfilePhotoUrl string    `json:"profile_photo_url"`
	BannerImageUrl  string    `json:"banner_image_url"`
	GalleryPhotoUrl []string  `json:"gallery_photo_url"`
	Restricted      bool      `json:"restricted" nullable:"true"`
	Verified        bool      `json:"verified" nullable:"true"`
	TotalFollowers  int       `json:"total_followers" default:"0"`
	GoogleMapUrl    string    `json:"google_maps_url"`
	IFrameMapUrl    string    `json:"iframe_maps_url"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`

	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
