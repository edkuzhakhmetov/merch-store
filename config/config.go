package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	if err := validateDatabase(); err != nil {
		return fmt.Errorf("database configuration error: %w", err)
	}

	if err := validateJWT(); err != nil {
		return fmt.Errorf("JWT configuration error: %w", err)
	}

	return nil
}

func validateDatabase() error {
	if os.Getenv("DB_HOST") == "" {
		return fmt.Errorf("database host address (DB_HOST) is not set")
	}

	if os.Getenv("DB_PORT") == "" {
		return fmt.Errorf("database port (DB_PORT) is not set")
	}

	if os.Getenv("DB_USER") == "" {
		return fmt.Errorf("database username (DB_USER) is not set")
	}

	return nil
}

func validateJWT() error {
	if os.Getenv("JWT_SECRET_KEY") == "" {
		return fmt.Errorf("JWT secret (JWT_SECRET_KEY) is not set")
	}

	return nil
}
