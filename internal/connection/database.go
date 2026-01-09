package connection

import (
	"fmt"
	"log"
	"time"

	"go-fiber-api/internal/config"
	"go-fiber-api/internal/features/auth"
	"go-fiber-api/internal/features/follow"
	"go-fiber-api/internal/features/inventory"
	"go-fiber-api/internal/features/merchant"
	"go-fiber-api/internal/features/products"
	"go-fiber-api/internal/features/transactions"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	cfg := config.Get()

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database connection string is empty")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := migrate(); err != nil {
		return nil, err
	}

	db = db.Debug()

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("database connected successfully")
	return db, nil
}

func migrate() error {
	log.Println("running database migrations...")

	dsn := config.Get().DatabaseURL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("failed to connect database for migration: %w", err)
	}

	if err := db.AutoMigrate(
		&auth.User{},
		&merchant.Merchant{},
		&products.Product{},
		&follow.Follow{},
		&transactions.Transaction{},
		&transactions.TransactionItem{},
		&inventory.StockMovement{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	log.Println("database migrated successfully")
	return nil
}
