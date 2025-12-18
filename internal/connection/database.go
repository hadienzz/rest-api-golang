package connection

import (
	"fmt"
	"log"

	"go-fiber-api/internal/config"
	"go-fiber-api/internal/features/auth"
	"go-fiber-api/internal/features/follow"
	"go-fiber-api/internal/features/merchant"
	"go-fiber-api/internal/features/products"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	cfg := config.Get()

	if cfg.Database == "" {
		return nil, fmt.Errorf("database connection string is empty")
	}

	db, err := gorm.Open(postgres.Open(cfg.Database), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	log.Println("database connected successfully")
	return db, nil
}

func migrate(db *gorm.DB) error {
	log.Println("running database migrations...")

	if err := db.AutoMigrate(
		&auth.User{},
		&merchant.Merchant{},
		&products.Product{},
		&follow.Follow{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	log.Println("database migrated successfully")
	return nil
}
