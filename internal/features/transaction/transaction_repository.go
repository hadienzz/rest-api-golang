package transaction

import (
	"errors"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(trx *Transaction) error
	FindByOrderID(orderID string) (*Transaction, error)
	FindByIdempotencyKey(key string) (*Transaction, error)
	UpdateStatusAndPaymentType(orderID string, status TransactionStatus, paymentType string) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(trx *Transaction) error {
	result := r.db.Create(trx)
	return result.Error
}

func (r *transactionRepository) FindByOrderID(orderID string) (*Transaction, error) {
	var trx Transaction
	result := r.db.Where("order_id = ?", orderID).First(&trx)
	if result.Error != nil {
		return nil, result.Error
	}
	return &trx, nil
}

func (r *transactionRepository) FindByIdempotencyKey(key string) (*Transaction, error) {
	var trx Transaction
	result := r.db.Where("idempotency_key = ?", key).First(&trx)
	if result.Error != nil {
		return nil, result.Error
	}
	return &trx, nil
}
func (r *transactionRepository) UpdateStatusAndPaymentType(
	orderID string,
	status TransactionStatus,
	paymentType string,
) error {

	updates := map[string]interface{}{
		"status":       status,
		"payment_type": paymentType,
	}

	result := r.db.
		Model(&Transaction{}).
		Where("order_id = ?", orderID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("transaction not found for order_id")
	}

	return nil
}
