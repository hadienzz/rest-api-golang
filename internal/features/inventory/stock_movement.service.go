package inventory

import "github.com/google/uuid"

type stockMovementService struct {
	repo StockMovementRepository
}

type StockMovementService interface {
	AddStockIn(productID uuid.UUID, quantity int) error
	AddStockOut(productID uuid.UUID, quantity int) error
	AddStockSale(productID uuid.UUID, quantity int) error
}

func NewStockMovementService(repo StockMovementRepository) StockMovementService {
	return &stockMovementService{repo: repo}
}

func (s *stockMovementService) AddStockIn(productID uuid.UUID, quantity int) error {
	return s.repo.AddStockIn(productID, quantity)
}

func (s *stockMovementService) AddStockOut(productID uuid.UUID, quantity int) error {
	return s.repo.AddStockOut(productID, quantity)
}

func (s *stockMovementService) AddStockSale(productID uuid.UUID, quantity int) error {
	return s.repo.AddStockSale(productID, quantity)
}
