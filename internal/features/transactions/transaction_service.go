package transactions

import (
	"errors"
	"fmt"

	"go-fiber-api/internal/config"
	"go-fiber-api/internal/features/inventory"
	"go-fiber-api/internal/features/products"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(userID uuid.UUID, req *CreateTransactionRequest) (*CreateTransactionResponse, error)
	HandleMidtransWebhook(req *MidtransNotificationRequest) error
	GetTransactionDetail(transactionID string) (*TransactionDetailResponse, error)
	GetTransactionsByUserID(userID uuid.UUID) ([]TransactionDetailResponse, error)
	ResumeTransactionByIdempotencyKey(idempotencyKey string) (*CreateTransactionResponse, error)
	GetTransactionsByMerchantID(merchantID uuid.UUID) ([]TransactionDTO, error)
}

type transactionService struct {
	db                    *gorm.DB
	transactionRepository TransactionRepository
	itemRepository        TransactionItemRepository
	productRepository     products.ProductRepository
	stockMovementRepo     inventory.StockMovementRepository
}

func NewTransactionService(
	db *gorm.DB,
	transactionRepo TransactionRepository,
	itemRepo TransactionItemRepository,
	productRepo products.ProductRepository,
	stockMovementRepo inventory.StockMovementRepository,
) TransactionService {
	return &transactionService{
		db:                    db,
		transactionRepository: transactionRepo,
		itemRepository:        itemRepo,
		productRepository:     productRepo,
		stockMovementRepo:     stockMovementRepo,
	}
}

func mapMidtransStatus(mtStatus string) TransactionStatus {
	switch mtStatus {
	case "settlement", "capture":
		return TransactionStatusPaid
	case "cancel", "expire", "deny":
		return TransactionStatusFailed
	default:
		return TransactionStatusPending
	}
}

func (s *transactionService) CreateTransaction(userID uuid.UUID, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	if req.IdempotencyKey == "" {
		return nil, fmt.Errorf("idempotency_key is required")
	}

	// Cek idempotency_key untuk mencegah duplikasi payment
	existing, err := s.transactionRepository.FindByIdempotencyKey(req.IdempotencyKey)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil && existing.ID != uuid.Nil {
		// Jika transaksi dengan idempotency_key ini sudah ada,
		// selalu kembalikan informasi transaksi yang sama
		if existing.SnapToken == "" || existing.RedirectURL == "" {
			return nil, fmt.Errorf("transaction already exists but missing payment info")
		}

		return &CreateTransactionResponse{
			OrderID:     existing.OrderID,
			SnapToken:   existing.SnapToken,
			RedirectURL: existing.RedirectURL,
			Status:      string(existing.Status),
		}, nil
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items cannot be empty")
	}

	productIDs := make([]uuid.UUID, 0, len(req.Items))
	for _, item := range req.Items {

		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than zero")
		}
		productIDs = append(productIDs, item.ProductID)
	}

	productsList, err := s.productRepository.GetProductsByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]products.Product, len(productsList))
	for _, p := range productsList {
		productMap[p.ID] = p
	}

	if len(productMap) != len(productIDs) {
		return nil, fmt.Errorf("one or more products not found")
	}

	orderID := fmt.Sprintf("ORDER-%s", uuid.NewString())
	totalAmount := decimal.NewFromInt(0)
	transactionItems := make([]TransactionItem, 0, len(req.Items))
	itemDetails := make([]midtrans.ItemDetails, 0, len(req.Items))

	for _, itemReq := range req.Items {
		product, ok := productMap[itemReq.ProductID]
		if !ok {
			return nil, fmt.Errorf("product not found")
		}

		if product.MerchantID != req.MerchantID {
			return nil, fmt.Errorf("product does not belong to merchant")
		}

		priceDecimal := product.Price
		if priceDecimal.LessThanOrEqual(decimal.NewFromInt(0)) {
			return nil, fmt.Errorf("invalid product price")
		}

		qtyDec := decimal.NewFromInt(int64(itemReq.Quantity))
		subtotal := priceDecimal.Mul(qtyDec)
		totalAmount = totalAmount.Add(subtotal)

		transactionItems = append(transactionItems, TransactionItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     priceDecimal,
			Subtotal:  subtotal,
		})

		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:   product.ID.String(),
			Name: product.Name,
			// Midtrans tetap butuh integer rupiah, pakai nilai dibulatkan
			Price: priceDecimal.Round(0).IntPart(),
			Qty:   int32(itemReq.Quantity),
		})
	}

	if totalAmount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return nil, fmt.Errorf("total amount must be greater than zero")
	}

	transaction := &Transaction{
		UserID:         userID,
		MerchantID:     req.MerchantID,
		OrderID:        orderID,
		Status:         TransactionStatusPending,
		TotalAmount:    totalAmount,
		IdempotencyKey: req.IdempotencyKey,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		trxRepo := NewTransactionRepository(tx)
		itemRepo := NewTransactionItemRepository(tx)

		if err := trxRepo.Create(transaction); err != nil {
			return err
		}

		for i := range transactionItems {
			transactionItems[i].TransactionID = transaction.ID
		}

		if err := itemRepo.BulkCreate(transactionItems); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	cfg := config.Get()
	var snapClient snap.Client
	snapClient.New(cfg.MidtransServerKey, midtrans.Sandbox)

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: transaction.OrderID,
			// GrossAmt harus integer rupiah untuk Midtrans
			GrossAmt: transaction.TotalAmount.Round(0).IntPart(),
		},
		Items: &itemDetails,
	}

	snapResp, err := snapClient.CreateTransaction(snapReq)

	// Library Midtrans kadang mengembalikan error interface namun HTTP 200
	// dan body berisi token. Jika snapResp tidak nil dan ada token,
	// kita anggap sukses dan abaikan err untuk menghindari bug error nil-pointer.
	if snapResp == nil {
		// Jika tidak ada response sama sekali, baru error kita propagasi.
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create snap transaction: empty response from Midtrans")
	}

	// Simpan token dan redirect URL ke dalam record transaksi
	updateErr := s.db.Model(&Transaction{}).
		Where("id = ?", transaction.ID).
		Updates(map[string]interface{}{
			"snap_token":   snapResp.Token,
			"redirect_url": snapResp.RedirectURL,
		}).Error

	if updateErr != nil {
		return nil, updateErr
	}

	transaction.SnapToken = snapResp.Token
	transaction.RedirectURL = snapResp.RedirectURL

	response := &CreateTransactionResponse{
		OrderID:     transaction.OrderID,
		SnapToken:   transaction.SnapToken,
		RedirectURL: transaction.RedirectURL,
		Status:      string(transaction.Status),
	}

	return response, nil
}

func (s *transactionService) HandleMidtransWebhook(
	req *MidtransNotificationRequest,
) error {

	if req.OrderID == "" {
		return errors.New("order_id is required")
	}

	transaction, err := s.transactionRepository.FindByOrderID(req.OrderID)
	if err != nil {
		return err
	}

	// ⛔ Jangan overwrite status final
	if transaction.Status == TransactionStatusPaid ||
		transaction.Status == TransactionStatusFailed {
		return nil
	}

	newStatus := mapMidtransStatus(req.TransactionStatus)

	return s.db.Transaction(func(dbTx *gorm.DB) error {
		trxRepo := s.transactionRepository.WithTx(dbTx)
		stockRepo := s.stockMovementRepo.WithTx(dbTx)

		if err := trxRepo.UpdateStatusAndPaymentType(
			req.OrderID,
			newStatus,
			req.PaymentType,
		); err != nil {
			return err
		}

		if newStatus == TransactionStatusPaid {
			for _, item := range transaction.Items {
				if err := stockRepo.AddStockSale(item.ProductID, item.Quantity); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (s *transactionService) GetTransactionDetail(id string) (*TransactionDetailResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("transaction_id is required")
	}

	tx, err := s.transactionRepository.GetTransactionsDetailByID(id)
	if err != nil {
		return nil, err
	}

	items := make([]TransactionItemResponse, 0, len(tx.Items))
	for _, item := range tx.Items {
		items = append(items, TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Subtotal:    item.Subtotal,
		})
	}

	resp := &TransactionDetailResponse{
		ID:             tx.ID,
		OrderID:        tx.OrderID,
		Status:         string(tx.Status),
		TotalAmount:    tx.TotalAmount,
		PaymentType:    tx.PaymentType,
		MerchantID:     tx.MerchantID,
		MerchantName:   tx.Merchant.Name,
		IdempotencyKey: tx.IdempotencyKey,
		CreatedAt:      tx.CreatedAt,
		Items:          items,
	}

	return resp, nil
}

func (s *transactionService) GetTransactionsByUserID(userID uuid.UUID) ([]TransactionDetailResponse, error) {
	var transactions []TransactionWithMerchant

	transactions, err := s.transactionRepository.GetTransactionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var resp []TransactionDetailResponse

	for _, tx := range transactions {
		resp = append(resp, TransactionDetailResponse{
			ID:             tx.ID,
			OrderID:        tx.OrderID,
			Status:         string(tx.Status),
			TotalAmount:    tx.TotalAmount,
			MerchantName:   tx.MerchantName,
			MerchantID:     tx.MerchantID,
			PaymentType:    tx.PaymentType,
			CreatedAt:      tx.CreatedAt,
			IdempotencyKey: tx.IdempotencyKey,
		})
	}

	return resp, nil

}

func (s *transactionService) ResumeTransactionByIdempotencyKey(idempotencyKey string) (*CreateTransactionResponse, error) {
	if idempotencyKey == "" {
		return nil, fmt.Errorf("idempotency_key is required")
	}

	tx, err := s.transactionRepository.FindByIdempotencyKey(idempotencyKey)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}

	if tx.SnapToken == "" || tx.RedirectURL == "" {
		return nil, fmt.Errorf("transaction exists but missing payment info")
	}

	return &CreateTransactionResponse{
		OrderID:     tx.OrderID,
		SnapToken:   tx.SnapToken,
		RedirectURL: tx.RedirectURL,
		Status:      string(tx.Status),
	}, nil
}

func (s *transactionService) GetTransactionsByMerchantID(merchantID uuid.UUID) ([]TransactionDTO, error) {

	transaction, err := s.transactionRepository.GetTransactionByMerchantID(merchantID)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
