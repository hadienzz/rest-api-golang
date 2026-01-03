package inventory

import "github.com/google/uuid"

type StockMovementDTO struct {
	ProductID uuid.UUID `json:"product_id" validate:"required,uuid4"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}
