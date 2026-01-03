package main

import (
	"go-fiber-api/internal/api"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/connection"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	db, err := connection.ConnectDB()
	err = config.InitSupabase()
	if err != nil {
		panic("failed to initialize Supabase or Database: " + err.Error())
	}

	config.InitMidtrans()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	api.RegisterAuthRoutes(app, db)
	api.RegisterMerchantRoutes(app, db)
	api.RegisterProductRoutes(app, db)
	api.RegisterFollowRoutes(app, db)
	api.RegisterTransactionRoutes(app, db)
	app.Listen(":8080")
}
