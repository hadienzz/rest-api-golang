package transaction

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionItemRepository interface {
	BulkCreate(items []TransactionItem) error
	FindByTransactionID(txID uuid.UUID) ([]TransactionItem, error)
}

type transactionItemRepository struct {
	db *gorm.DB
}

func NewTransactionItemRepository(db *gorm.DB) TransactionItemRepository {
	return &transactionItemRepository{db: db}
}

func (r *transactionItemRepository) BulkCreate(items []TransactionItem) error {
	if len(items) == 0 {
		return nil
	}
	result := r.db.Create(&items)
	return result.Error
}

func (r *transactionItemRepository) FindByTransactionID(txID uuid.UUID) ([]TransactionItem, error) {
	var items []TransactionItem
	result := r.db.Where("transaction_id = ?", txID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}
