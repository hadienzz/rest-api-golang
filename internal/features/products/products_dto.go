package products

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	MerchantID      uuid.UUID `json:"-" form:"-"`
	Name            string    `json:"name" form:"name" validate:"required"`
	Description     string    `json:"description" form:"description"`
	Price           string    `json:"price" form:"price" validate:"required"`
	Quantity        int       `json:"quantity" form:"quantity" validate:"required"`
	ProductPhotoUrl string    `json:"product_photo_url" form:"product_photo_url"`
}

type ProductDTO struct {
	ID              uuid.UUID       `json:"id"`
	MerchantID      uuid.UUID       `json:"merchant_id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Price           decimal.Decimal `json:"price"`
	Quantity        int             `json:"quantity"`
	ProductPhotoUrl string          `json:"product_photo_url"`
	CreatedAt       sql.NullTime    `json:"created_at"`
	UpdatedAt       sql.NullTime    `json:"updated_at"`
}

type BulkDeleteProductRequest struct {
	ProductIDs []uuid.UUID `json:"product_ids"`
}
