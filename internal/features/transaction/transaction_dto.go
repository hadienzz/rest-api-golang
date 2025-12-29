package transaction

import "github.com/google/uuid"

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
