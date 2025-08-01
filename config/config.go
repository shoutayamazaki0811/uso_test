package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	JWTSecret           string
	Port                string
	StripeSecretKey     string
	StripeWebhookSecret string
	ResendAPIKey        string
	AdminPassword       string
	BaseURL             string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://localhost/uso_db?sslmode=disable"),
		JWTSecret:           getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		Port:                getEnv("PORT", "8080"),
		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		ResendAPIKey:        getEnv("RESEND_API_KEY", ""),
		AdminPassword:       getEnv("ADMIN_PASSWORD", "admin123"),
		BaseURL:             getEnv("BASE_URL", "http://localhost:8080"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}