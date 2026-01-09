package config

import (
	"os"

	"github.com/joho/godotenv"
)

func Get() *Config {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		DatabaseDirectURL: os.Getenv("DATABASE_DIRECT_URL"),
		DatabasePoolerURL: os.Getenv("DATABASE_POOLER_URL"),
		Port:              os.Getenv("PORT"),
		JwtKey:            os.Getenv("JWT_SECRET_KEY"),
		SupabaseURL:       os.Getenv("SUPABASE_URL"),

		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		MidtransServerKey:  os.Getenv("MIDTRANS_SERVER_KEY"),
		MidtransMerchantID: os.Getenv("MIDTRANS_MERCHANT_ID"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
	}
}
