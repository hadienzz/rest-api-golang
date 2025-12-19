package merchant

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Merchant struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID          uuid.UUID      `gorm:"column:user_id;not null;index"`
	Name            string         `gorm:"column:name;type:varchar(100);not null"`
	Description     string         `gorm:"column:description;type:text"`
	Type            string         `gorm:"column:type;type:varchar(50)"`
	Location        string         `gorm:"column:location"`
	ProfilePhotoUrl string         `gorm:"column:profile_photo_url;type:text"`
	BannerImageUrl  string         `gorm:"column:banner_image_url;type:text"`
	GalleryPhotoUrl pq.StringArray `gorm:"column:gallery_photo_url;type:text[]"`
	Restricted      bool           `gorm:"column:restricted;default:false"`
	Verified        bool           `gorm:"column:verified;default:false"`
	TotalFollowers  int            `gorm:"column:total_followers;default:0"`
	GoogleMapUrl    string         `gorm:"column:google_maps_url;type:text"`
	IFrameMapUrl    string         `gorm:"column:iframe_maps_url;type:text"`
	Latitude        float64        `gorm:"type:decimal(9,6);not null"`
	Longitude       float64        `gorm:"type:decimal(9,6);not null"`

	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
