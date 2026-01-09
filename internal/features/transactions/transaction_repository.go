package transactions

import (
	"errors"
	"go-fiber-api/internal/util/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(trx *Transaction) error
	FindByOrderID(orderID string) (*Transaction, error)
	FindByIdempotencyKey(key string) (*Transaction, error)
	UpdateStatusAndPaymentType(orderID string, status TransactionStatus, paymentType string) error
	GetTransactionsByUserID(userID uuid.UUID) ([]TransactionWithMerchant, error)
	GetTransactionsDetailByID(orderID string) (*Transaction, error)
	GetTransactionByMerchantID(MerchantID uuid.UUID, page int, limit int) ([]TransactionDTO, int64, error)
	WithTx(tx *gorm.DB) *transactionRepository
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) WithTx(tx *gorm.DB) *transactionRepository {
	return &transactionRepository{db: tx}
}

func (r *transactionRepository) Create(trx *Transaction) error {
	result := r.db.Create(trx)
	return result.Error
}

func (r *transactionRepository) FindByOrderID(orderID string) (*Transaction, error) {
	var trx Transaction
	result := r.db.
		Preload("Items").
		Where("order_id = ?", orderID).
		First(&trx)
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

func (r *transactionRepository) GetTransactionsByUserID(userID uuid.UUID) ([]TransactionWithMerchant, error) {
	var transactions []TransactionWithMerchant

	err := r.db.
		Table("transactions").
		Select(`
		transactions.*,
		merchants.name AS merchant_name
	`).
		Joins("JOIN merchants ON merchants.id = transactions.merchant_id").
		Where("transactions.user_id = ?", userID).
		Find(&transactions).
		Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetTransactionsDetailByID(transactionID string) (*Transaction, error) {
	var transaction Transaction

	err := r.db.
		Preload("Items.Product").
		Preload("Merchant", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Where("id = ?", transactionID).
		First(&transaction).
		Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) GetTransactionByMerchantID(MerchantID uuid.UUID, page int, limit int) ([]TransactionDTO, int64, error) {
	var transactions []TransactionDTO

	var totalItems int64
	offset := pagination.GetOffset(page, limit)

	if err := r.db.Model(&Transaction{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Table("transactions").
		Where("merchant_id = ?", MerchantID).
		Limit(10).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, totalItems, nil
}
