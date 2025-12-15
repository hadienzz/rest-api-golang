package products

type ProductService interface {
	CreateProduct(product *Product) (*Product, error)
	GetMerchantProducts(merchantID string) ([]ProductDTO, error)
}

type productService struct {
	productRepository ProductRepository
}

func NewProductService(productRepository ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (ps *productService) CreateProduct(product *Product) (*Product, error) {
	return ps.productRepository.CreateProduct(product)
}

func (ps *productService) GetMerchantProducts(merchantID string) ([]ProductDTO, error) {
	products, err := ps.productRepository.GetMerchantProducts(merchantID)
	if err != nil {
		return nil, err
	}

	responses := make([]ProductDTO, 0, len(products))

	for _, e := range products {
		responses = append(responses, ProductDTO{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			Price:       e.Price,
			Stock:       e.Stock,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		})
	}

	return responses, nil
}
