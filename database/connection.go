package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shodan-backend/models"
)

// InitDB establishes the database connection and runs migrations.
// It expects DATABASE_URL to be set (optionally loaded from a local .env file).
func InitDB() (*gorm.DB, error) {
	// Load .env if present; ignore error if file doesn't exist
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := db.AutoMigrate(&models.Result{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
