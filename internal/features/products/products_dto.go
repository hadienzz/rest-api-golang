package products

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	MerchantID  uuid.UUID       `json:"merchant_id" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" validate:"required,min=0"`
	Quantity    int             `json:"quantity" validate:"required,min=0"`
}

type ProductDTO struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`

	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
