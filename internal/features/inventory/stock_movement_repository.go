package inventory

import (
	"errors"
	"go-fiber-api/internal/features/products"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockMovementRepository interface {
	WithTx(tx *gorm.DB) StockMovementRepository
	AddStockIn(productID uuid.UUID, quantity int) error
	AddStockOut(productID uuid.UUID, quantity int) error
	AddStockSale(productID uuid.UUID, quantity int) error
}

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{
		db: db,
	}
}

func (r *stockMovementRepository) WithTx(tx *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: tx}
}

func (r *stockMovementRepository) AddStockIn(productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	movement := &StockMovement{
		ProductID: productID,
		Quantity:  quantity,
		Type:      StockIn,
	}

	if err := r.db.Create(movement).Error; err != nil {
		return err
	}

	result := r.db.Model(&products.Product{}).
		Where("id = ?", productID).
		Update("quantity", gorm.Expr("quantity + ?", quantity))

	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}

	return result.Error
}

func (r *stockMovementRepository) AddStockOut(productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	result := r.db.Model(&products.Product{}).
		Where("id = ? AND quantity >= ?", productID, quantity).
		Update("quantity", gorm.Expr("quantity - ?", quantity))

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock or product not found")
	}

	if result.Error != nil {
		return result.Error
	}

	movement := &StockMovement{
		ProductID: productID,
		Type:      StockOut,
		Quantity:  quantity,
	}

	return r.db.Create(movement).Error
}

func (r *stockMovementRepository) AddStockSale(productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	result := r.db.Model(&products.Product{}).
		Where("id = ? AND quantity >= ?", productID, quantity).
		Update("quantity", gorm.Expr("quantity - ?", quantity))

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock or product not found")
	}

	if result.Error != nil {
		return result.Error
	}

	movement := &StockMovement{
		ProductID: productID,
		Type:      StockSale,
		Quantity:  quantity,
	}

	return r.db.Create(movement).Error
}
