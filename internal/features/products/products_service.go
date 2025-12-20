package products

// import "go-fiber-api/internal/features/merchant"

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductService interface {
	CreateProduct(product *CreateProductRequest) (*ProductDTO, error)
	GetMerchantProducts(merchantID uuid.UUID) ([]ProductDTO, error)
	DeleteMerchantProduct(productID []uuid.UUID, merchantID uuid.UUID) error
}

type productService struct {
	productRepository ProductRepository
	merchantAdapter   MerchantServiceContract
}

func NewProductService(productRepository ProductRepository, merchantAdapter MerchantServiceContract) ProductService {
	return &productService{
		productRepository: productRepository,
		merchantAdapter:   merchantAdapter,
	}
}

func (ps *productService) CreateProduct(req *CreateProductRequest) (*ProductDTO, error) {
	priceDecimal, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, err
	}

	product := &Product{
		MerchantID:      req.MerchantID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           priceDecimal,
		Quantity:        req.Quantity,
		ProductPhotoUrl: req.ProductPhotoUrl,
	}

	createdProduct, err := ps.productRepository.CreateProduct(product)

	if err != nil {
		return nil, err
	}

	return &ProductDTO{
		ID:              createdProduct.ID,
		MerchantID:      createdProduct.MerchantID,
		Name:            createdProduct.Name,
		Description:     createdProduct.Description,
		Price:           createdProduct.Price,
		Quantity:        createdProduct.Quantity,
		ProductPhotoUrl: createdProduct.ProductPhotoUrl,
		CreatedAt:       createdProduct.CreatedAt,
		UpdatedAt:       createdProduct.UpdatedAt,
	}, nil
}

func (ps *productService) GetMerchantProducts(merchantID uuid.UUID) ([]ProductDTO, error) {
	products, err := ps.productRepository.GetMerchantProducts(merchantID)

	if err != nil {
		return nil, err
	}

	responses := make([]ProductDTO, 0, len(products))

	for _, e := range products {
		responses = append(responses, ProductDTO{
			ID:              e.ID,
			MerchantID:      e.MerchantID,
			Name:            e.Name,
			Description:     e.Description,
			Price:           e.Price,
			Quantity:        e.Quantity,
			ProductPhotoUrl: e.ProductPhotoUrl,
			CreatedAt:       e.CreatedAt,
			UpdatedAt:       e.UpdatedAt,
		})
	}

	return responses, nil
}

func (ps *productService) DeleteMerchantProduct(productID []uuid.UUID, merchantID uuid.UUID) error {
	return ps.productRepository.DeleteMerchantProduct(productID, merchantID)
}
