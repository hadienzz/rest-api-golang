package products

// import "go-fiber-api/internal/features/merchant"

type ProductService interface {
	CreateProduct(product *CreateProductRequest) (*ProductDTO, error)
	GetMerchantProducts(merchantID string) ([]ProductDTO, error)
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

	product := &Product{
		MerchantID:  req.MerchantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}

	createdProduct, err := ps.productRepository.CreateProduct(product)

	if err != nil {
		return nil, err
	}

	return &ProductDTO{
		ID:          createdProduct.ID,
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
		Quantity:    createdProduct.Quantity,
		CreatedAt:   createdProduct.CreatedAt,
		UpdatedAt:   createdProduct.UpdatedAt,
	}, nil
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
			Quantity:    e.Quantity,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
		})
	}

	return responses, nil
}
