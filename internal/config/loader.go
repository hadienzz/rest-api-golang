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
		Database:           os.Getenv("DATABASE_URL"),
		Port:               os.Getenv("PORT"),
		JwtKey:             os.Getenv("JWT_SECRET_KEY"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
	}
}
