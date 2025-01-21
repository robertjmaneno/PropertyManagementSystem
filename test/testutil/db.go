package testutil

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewTestDB creates a new test PostgreSQL database
func NewTestDB() (*gorm.DB, error) {
	// Get test database configuration from environment or use defaults
	host := getEnv("TEST_DB_HOST", "localhost")
	port := getEnv("TEST_DB_PORT", "5432")
	user := getEnv("TEST_DB_USER", "postgres")
	password := getEnv("TEST_DB_PASSWORD", "postgres")
	dbname := getEnv("TEST_DB_NAME", "template_test")
	schema := getEnv("TEST_DB_SCHEMA", "test")

	// Build DSN with schema
	query := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		host, port, user, password, dbname, schema)

	db, err := gorm.Open(postgres.Open(query), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create schema if it doesn't exist
	if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error; err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	// Set search path
	if err := db.Exec(fmt.Sprintf("SET search_path TO %s,public", schema)).Error; err != nil {
		return nil, fmt.Errorf("failed to set search path: %w", err)
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
