package connection

import (
	"fmt"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/features/auth"
	"go-fiber-api/internal/features/follow"
	"go-fiber-api/internal/features/merchant"
	"go-fiber-api/internal/features/products"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var databaseInstance *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	loadedConfig := config.Get()

	databaseUrl := loadedConfig.Database

	if databaseUrl == "" {
		panic("Database connection string is empty")
	}

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate di sini
	if err := migrate(db); err != nil {
		return nil, err
	}

	log.Print("Database is connected")
	databaseInstance = db
	return db, nil
}

func migrate(db *gorm.DB) error {
	log.Println("running database migrations...")

	// daftar semua entity di sini
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

func GetDatabase() *gorm.DB {
	if databaseInstance == nil {
		panic("Database is not connected")
	}

	return databaseInstance
}
