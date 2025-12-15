package follow

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type FollowerDTO struct {
	ID uuid.UUID `json:"id"`

	UserID     uuid.UUID    `json:"-"`
	MerchantID uuid.UUID    `json:"merchant_id"`
	CreatedAt  sql.NullTime `json:"created_at"`
}

type FollowRequest struct {
	UserID     uuid.UUID `json:"-"`
	MerchantID uuid.UUID `json:"merchant_id"`
}

type FollowResponse struct {
	IsFollowing bool      `json:"is_following"`
	MerchantID  uuid.UUID `json:"merchant_id"`
	FollowedAt  time.Time `json:"followed_at" nullable:"true"`
}
