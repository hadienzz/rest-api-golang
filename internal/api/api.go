package api

import (
	"go-fiber-api/internal/features/auth"
	"go-fiber-api/internal/features/follow"
	"go-fiber-api/internal/features/merchant"
	"go-fiber-api/internal/features/products"
	"go-fiber-api/internal/features/transaction"

	// "go-fiber-api/internal/features/products"
	"go-fiber-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/auth")
	userRepository := auth.NewUserRepository(db)
	authService := auth.NewAuthService(userRepository)
	authHandler := auth.NewHandler(authService)

	api.Post("/logout", authHandler.LogoutUser)
	api.Post("/register", authHandler.RegisterUser)
	api.Post("/login", authHandler.LoginUser)

	api.Get("/user", middleware.AuthRequired, authHandler.GetUser)
}

func RegisterMerchantRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/merchant")
	merchantRepo := merchant.NewMerchantRepository(db)
	merchantService := merchant.NewMerchantService(merchantRepo)
	merchantHandler := merchant.NewMerchantHandler(merchantService)

	api.Post("/create", middleware.AuthRequired, merchantHandler.AddMerchant)

	api.Get("/all", merchantHandler.GetAllMerchant)
	api.Get("/my-summary", middleware.AuthRequired, merchantHandler.GetMyMerchantsSummary)
	api.Get("/my-merchant/:id", middleware.AuthRequired, merchantHandler.GetMyMerchantDashboard)
	api.Get("/display", merchantHandler.GetMerchantDisplay)
	api.Get("/:id", merchantHandler.GetMerchantById)
}

func RegisterProductRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/products")

	productRepo := products.NewProductRepository(db)
	merchantRepo := merchant.NewMerchantRepository(db)
	merchantService := merchant.NewMerchantService(merchantRepo)
	merchantAdapter := merchant.NewMerchantServiceAdapter(merchantService)

	productService := products.NewProductService(productRepo, merchantAdapter)
	productHandler := products.NewProductHandler(productService, merchantAdapter)

	api.Get("/dashboard/:merchant_id", middleware.AuthRequired, productHandler.GetMerchantProductsDashboard)
	api.Get("/merchant/:id", productHandler.GetMerchantProducts)
	api.Post("/bulk-delete", middleware.AuthRequired, productHandler.BulkDeleteMerchantProducts)
	api.Post("/add/:merchant_id", middleware.AuthRequired, productHandler.CreateProduct)
	// api.Get("/me")
}

func RegisterFollowRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/follow")
	followRepo := follow.NewFollowersRepository(db)
	followService := follow.NewFollowService(followRepo)
	followHandler := follow.NewFollowController(followService)

	api.Post("/merchant/:id", middleware.AuthRequired, followHandler.FollowMerchant)
	api.Delete("/merchant/:id", middleware.AuthRequired, followHandler.UnfollowMerchant)
	api.Get("/merchant/:id/status", middleware.AuthRequired, followHandler.GetMerchantFollowStatus)
	// api.Get("/merchant", middleware.AuthRequired, follow)
}

func RegisterTransactionRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/transactions")

	transactionRepo := transaction.NewTransactionRepository(db)
	transactionItemRepo := transaction.NewTransactionItemRepository(db)
	productRepo := products.NewProductRepository(db)

	transactionService := transaction.NewTransactionService(db, transactionRepo, transactionItemRepo, productRepo)
	transactionHandler := transaction.NewTransactionHandler(transactionService)

	api.Post("/", middleware.AuthRequired, transactionHandler.CreateTransaction)
	api.Post("/webhook/midtrans", transactionHandler.HandleMidtransWebhook)
}
