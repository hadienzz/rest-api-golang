package products

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	MerchantID  string          `json:"-"`
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price" validate:"required,min=0"`
	Stock       int             `json:"stock" validate:"required,min=0"`
}

type ProductDTO struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Stock       int             `json:"stock"`

	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
