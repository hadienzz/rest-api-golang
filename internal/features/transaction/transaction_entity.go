package transaction

import (
	"time"

	"go-fiber-api/internal/features/auth"
	"go-fiber-api/internal/features/merchant"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusPaid    TransactionStatus = "PAID"
	TransactionStatusFailed  TransactionStatus = "FAILED"
)

type Transaction struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	MerchantID uuid.UUID `gorm:"type:uuid;not null;index"`

	// IdempotencyKey digunakan untuk mencegah duplikasi payment/order
	IdempotencyKey string `gorm:"type:varchar(100);not null;uniqueIndex"`

	OrderID     string            `gorm:"type:varchar(100);not null;uniqueIndex"`
	Status      TransactionStatus `gorm:"type:varchar(50);not null;default:'PENDING'"`
	TotalAmount int64             `gorm:"type:bigint;not null"`
	PaymentType string            `gorm:"type:varchar(50)"`
	SnapToken   string            `gorm:"type:text"`
	RedirectURL string            `gorm:"type:text"`

	// Relations
	User     auth.User         `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Merchant merchant.Merchant `gorm:"foreignKey:MerchantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Items    []TransactionItem `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
