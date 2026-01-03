package transactions

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateTransactionItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type CreateTransactionRequest struct {
	MerchantID     uuid.UUID                      `json:"merchant_id"`
	Items          []CreateTransactionItemRequest `json:"items"`
	IdempotencyKey string                         `json:"idempotency_key"`
}

type CreateTransactionResponse struct {
	OrderID     string `json:"order_id"`
	SnapToken   string `json:"snap_token"`
	RedirectURL string `json:"redirect_url"`
	Status      string `json:"status"`
}

// MidtransNotificationRequest mewakili payload penting dari webhook Midtrans
type MidtransNotificationRequest struct {
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
}

// Response untuk detail transaksi beserta item-nya

type TransactionItemResponse struct {
	ID          uuid.UUID       `json:"id"`
	ProductID   uuid.UUID       `json:"product_id"`
	ProductName string          `json:"product_name"`
	Quantity    int             `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	Subtotal    decimal.Decimal `json:"subtotal"`
}

type TransactionDetailResponse struct {
	ID             uuid.UUID                 `json:"id"`
	OrderID        string                    `json:"order_id"`
	Status         string                    `json:"status"`
	TotalAmount    decimal.Decimal           `json:"total_amount"`
	PaymentType    string                    `json:"payment_type"`
	MerchantID     uuid.UUID                 `json:"merchant_id"`
	MerchantName   string                    `json:"merchant_name"`
	IdempotencyKey string                    `json:"idempotency_key"`
	CreatedAt      time.Time                 `json:"created_at"`
	Items          []TransactionItemResponse `json:"items"`
}

type TransactionDTO struct {
	ID          uuid.UUID       `json:"id"`
	OrderID     string          `json:"order_id"`
	Status      string          `json:"status"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	PaymentType string          `json:"payment_type"`
	MerchantID  uuid.UUID       `json:"merchant_id"`

	CreatedAt time.Time `json:"created_at"`
}

type TransactionWithMerchant struct {
	Transaction
	MerchantName string `json:"merchant_name" gorm:"column:merchant_name"`
}
