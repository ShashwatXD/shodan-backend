package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shodan-backend/models"
)

// InitDB initializes the PostgreSQL database connection and runs migrations.
func InitDB() (*gorm.DB, error) {
	// PostgreSQL connection string
	// You can also use environment variables for better security
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://shodan_user:RR3VRkIbn3DF34A0PDFI5IeEu0GqjSUR@dpg-d49322je5dus73ch5ej0-a.oregon-postgres.render.com/shodan?sslmode=require"
	}

	fmt.Printf("Attempting to connect to database with DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Auto-migrate the Result model
	if err := db.AutoMigrate(&models.Result{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
