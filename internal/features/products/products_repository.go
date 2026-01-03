package products

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	WithTx(tx *gorm.DB) ProductRepository

	CreateProduct(product *Product) (*Product, error)
	FindByUserID(userID string) ([]Product, error)
	GetMerchantProducts(merchantID uuid.UUID) ([]Product, error)
	DeleteMerchantProduct(productID []uuid.UUID, merchantID uuid.UUID) error
	GetMerchantProductsDashboard(merchantID uuid.UUID) ([]Product, error)
	GetProductsByIDs(ids []uuid.UUID) ([]Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (pr *productRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepository{
		db: tx,
	}
}

func (pr *productRepository) CreateProduct(product *Product) (*Product, error) {
	result := pr.db.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (pr *productRepository) FindByUserID(userID string) ([]Product, error) {
	var products []Product

	// Cari produk berdasarkan user_id pemilik merchant lewat join
	result := pr.db.
		Joins("JOIN merchants ON merchants.id = products.merchant_id").
		Where("merchants.user_id = ?", userID).
		Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (pr *productRepository) GetMerchantProducts(merchantID uuid.UUID) ([]Product, error) {
	var products []Product
	result := pr.db.Where("merchant_id = ?", merchantID).Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (pr *productRepository) DeleteMerchantProduct(productID []uuid.UUID, merchantID uuid.UUID) error {
	result := pr.db.Where("id IN ?", productID).Delete(&Product{})
	log.Printf("delete products | rows=%d | err=%v", result.RowsAffected, result.Error)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (pr *productRepository) GetMerchantProductsDashboard(merchantID uuid.UUID) ([]Product, error) {
	var products []Product

	result := pr.db.Where("merchant_id = ?", merchantID).Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (pr *productRepository) GetProductsByIDs(ids []uuid.UUID) ([]Product, error) {
	var products []Product
	if len(ids) == 0 {
		return products, nil
	}

	result := pr.db.Where("id IN ?", ids).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}
