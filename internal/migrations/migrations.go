package migrations

import (
	"fmt"
	"time"

	"github.com/yourusername/projectname/internal/config"
	"github.com/yourusername/projectname/internal/domain"
	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"uniqueIndex;not null;type:varchar(255)"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// DatabaseDialect represents different database types
type DatabaseDialect interface {
	CreateSchema(db *gorm.DB, schema string) error
	SetSearchPath(db *gorm.DB, schema string) error
	UpdateNullFields(db *gorm.DB, table string, fields map[string]interface{}, conditions string) error
}

// PostgresDialect implements DatabaseDialect for PostgreSQL
type PostgresDialect struct{}

func (d PostgresDialect) CreateSchema(db *gorm.DB, schema string) error {
	if schema == "" || schema == "public" {
		return nil
	}
	return db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error
}

func (d PostgresDialect) SetSearchPath(db *gorm.DB, schema string) error {
	if schema == "" || schema == "public" {
		return nil
	}
	return db.Exec(fmt.Sprintf("SET search_path TO %s,public", schema)).Error
}

func (d PostgresDialect) UpdateNullFields(db *gorm.DB, table string, fields map[string]interface{}, conditions string) error {
	query := fmt.Sprintf("UPDATE %s SET ", table)
	updates := make([]string, 0, len(fields))
	for field, value := range fields {
		updates = append(updates, fmt.Sprintf("%s = COALESCE(%s, %v)", field, field, value))
	}
	query += fmt.Sprintf("%s WHERE %s", updates[0], conditions)
	return db.Exec(query).Error
}

// MySQLDialect implements DatabaseDialect for MySQL
type MySQLDialect struct{}

func (d MySQLDialect) CreateSchema(db *gorm.DB, schema string) error {
	if schema == "" {
		return nil
	}
	return db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", schema)).Error
}

func (d MySQLDialect) SetSearchPath(db *gorm.DB, schema string) error {
	if schema == "" {
		return nil
	}
	return db.Exec(fmt.Sprintf("USE %s", schema)).Error
}

func (d MySQLDialect) UpdateNullFields(db *gorm.DB, table string, fields map[string]interface{}, conditions string) error {
	query := fmt.Sprintf("UPDATE %s SET ", table)
	updates := make([]string, 0, len(fields))
	for field, value := range fields {
		updates = append(updates, fmt.Sprintf("%s = IFNULL(%s, %v)", field, field, value))
	}
	query += fmt.Sprintf("%s WHERE %s", updates[0], conditions)
	return db.Exec(query).Error
}

// SQLiteDialect implements DatabaseDialect for SQLite
type SQLiteDialect struct{}

func (d SQLiteDialect) CreateSchema(db *gorm.DB, schema string) error {
	// SQLite doesn't support schemas, so we do nothing
	return nil
}

func (d SQLiteDialect) SetSearchPath(db *gorm.DB, schema string) error {
	// SQLite doesn't support schemas, so we do nothing
	return nil
}

func (d SQLiteDialect) UpdateNullFields(db *gorm.DB, table string, fields map[string]interface{}, conditions string) error {
	query := fmt.Sprintf("UPDATE %s SET ", table)
	updates := make([]string, 0, len(fields))
	for field, value := range fields {
		updates = append(updates, fmt.Sprintf("%s = COALESCE(%s, %v)", field, field, value))
	}
	query += fmt.Sprintf("%s WHERE %s", updates[0], conditions)
	return db.Exec(query).Error
}

// getDialect returns the appropriate dialect based on the driver
func getDialect(driver string) DatabaseDialect {
	switch driver {
	case "postgres":
		return PostgresDialect{}
	case "mysql":
		return MySQLDialect{}
	case "sqlite3":
		return SQLiteDialect{}
	default:
		return PostgresDialect{} // Default to PostgreSQL
	}
}

// Migrate runs all migrations
func Migrate(db *gorm.DB) error {
	// Load config to get schema and driver
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get the appropriate dialect
	dialect := getDialect(cfg.Database.Driver)

	// Ensure schema exists
	if err := dialect.CreateSchema(db, cfg.Database.Schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Set search path/database
	if err := dialect.SetSearchPath(db, cfg.Database.Schema); err != nil {
		return fmt.Errorf("failed to set schema: %w", err)
	}

	// Create migrations table if it doesn't exist using raw SQL
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	if err := db.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Register migrations here
	migrations := []struct {
		name string
		fn   func(*gorm.DB, DatabaseDialect) error
	}{
		{"create_users_table", createUsersTable},
		{"create_communities_table", createCommunitiesTable},
		{"create_buildings_table", createBuildingsTable},
		{"create_units_table", createUnitsTable},
		// Add new migrations here
	}

	// Run migrations
	for _, m := range migrations {
		var count int64
		if err := db.Model(&Migration{}).Where("name = ?", m.name).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count > 0 {
			// Migration already executed
			continue
		}

		// Run migration with dialect
		if err := m.fn(db, dialect); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", m.name, err)
		}

		// Record migration
		if err := db.Exec("INSERT INTO migrations (name) VALUES (?)", m.name).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", m.name, err)
		}
	}

	return nil
}

// Rollback rolls back the last migration
func Rollback(db *gorm.DB) error {
	// Load config to get schema and driver
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get the appropriate dialect
	dialect := getDialect(cfg.Database.Driver)

	// Set search path/database
	if err := dialect.SetSearchPath(db, cfg.Database.Schema); err != nil {
		return fmt.Errorf("failed to set schema: %w", err)
	}

	// Get the last migration
	var lastMigration struct {
		Name string
	}
	if err := db.Raw("SELECT name FROM migrations ORDER BY id DESC LIMIT 1").Scan(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}
	// Register rollbacks here
	rollbacks := map[string]func(*gorm.DB, DatabaseDialect) error{
		"create_users_table":       rollbackUsersTable,
		"create_communities_table": rollbackCommunitiesTable,
		"create_buildings_table":   rollbackBuildingsTable,
		"create_units_table":       rollbackUnitsTable,
		// Add new rollbacks here
	}
	// Run rollback
	if fn, ok := rollbacks[lastMigration.Name]; ok {
		if err := fn(db, dialect); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", lastMigration.Name, err)
		}

		// Remove migration record
		if err := db.Exec("DELETE FROM migrations WHERE name = ?", lastMigration.Name).Error; err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", lastMigration.Name, err)
		}
	}

	return nil
}

// Migration: Create users table
func createUsersTable(db *gorm.DB, dialect DatabaseDialect) error {
	// Create the table with all required fields
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Check if we need to set default values
	var count int64
	if err := db.Model(&domain.User{}).Where("organization_id IS NULL OR branch_id IS NULL").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for null values: %w", err)
	}

	// Only update if there are records needing defaults
	if count > 0 {
		defaultFields := map[string]interface{}{
			"organization_id": 1,
			"branch_id":       "'default'",
		}
		if err := dialect.UpdateNullFields(db, "users", defaultFields, "organization_id IS NULL OR branch_id IS NULL"); err != nil {
			return fmt.Errorf("failed to set default values: %w", err)
		}
	}

	return nil
}

// Rollback: Drop users table
func rollbackUsersTable(db *gorm.DB, dialect DatabaseDialect) error {
	return db.Migrator().DropTable(&domain.User{})
}

// Migration: Create buildings table
func createBuildingsTable(db *gorm.DB, dialect DatabaseDialect) error {
	// Create the table for the Building model
	if err := db.AutoMigrate(&domain.Building{}); err != nil {
		return fmt.Errorf("failed to create buildings table: %w", err)
	}

	// Check if we need to set default values (if any)
	var count int64
	if err := db.Model(&domain.Building{}).Where("description IS NULL").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for null values building table: %w", err)
	}

	// Only update if there are records needing defaults
	if count > 0 {
		defaultFields := map[string]interface{}{
			"description": "'No description available'",
		}
		if err := dialect.UpdateNullFields(db, "buildings", defaultFields, "description IS NULL"); err != nil {
			return fmt.Errorf("failed to set default values for buildings: %w", err)
		}
	}

	// Only update if there are records needing defaults
	if count > 0 {
		defaultFields := map[string]interface{}{
			"organization_id": 1,
			"branch_id":       "'default'",
		}
		if err := dialect.UpdateNullFields(db, "buildings", defaultFields, "organization_id IS NULL OR branch_id IS NULL"); err != nil {
			return fmt.Errorf("failed to set default values: %w", err)
		}
	}
	return nil
}

// Rollback: Drop buildings table
func rollbackBuildingsTable(db *gorm.DB, dialect DatabaseDialect) error {
	// Drop the buildings table along with constraints
	if err := db.Migrator().DropTable(&domain.Building{}); err != nil {
		return fmt.Errorf("failed to rollback buildings table: %w", err)
	}
	return nil
}

// Migration: Create communities table
func createCommunitiesTable(db *gorm.DB, dialect DatabaseDialect) error {
	// Create the table for the Community model
	if err := db.AutoMigrate(&domain.Community{}); err != nil {
		return fmt.Errorf("failed to create communities table: %w", err)
	}

	// Check if we need to set default values (if any)
	var count int64
	if err := db.Model(&domain.Community{}).Where("description IS NULL").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for null values in communities table: %w", err)
	}

	// Only update if there are records needing defaults
	if count > 0 {
		defaultFields := map[string]interface{}{
			"description": "'No description available'",
		}
		if err := dialect.UpdateNullFields(db, "communities", defaultFields, "description IS NULL"); err != nil {
			return fmt.Errorf("failed to set default values for communities: %w", err)
		}
	}

	// Only update if there are records needing defaults
	if count > 0 {
		defaultFields := map[string]interface{}{
			"organization_id": 1,
			"branch_id":       "'default'",
		}
		if err := dialect.UpdateNullFields(db, "buildings", defaultFields, "organization_id IS NULL OR branch_id IS NULL"); err != nil {
			return fmt.Errorf("failed to set default values: %w", err)
		}
	}
	return nil
}

// Rollback: Drop communities table
func rollbackCommunitiesTable(db *gorm.DB, dialect DatabaseDialect) error {
	return db.Migrator().DropTable(&domain.Community{})
}

// Migration: Create units table
func createUnitsTable(db *gorm.DB, dialect DatabaseDialect) error {
	// Create the table for the Unit model
	if err := db.AutoMigrate(&domain.Unit{}); err != nil {
		return fmt.Errorf("failed to create units table: %w", err)
	}

	// Check if we need to set default values (if any)
	var count int64
	if err := db.Model(&domain.Unit{}).Where("description IS NULL").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check for null values in units table: %w", err)
	}

	// Only update if there are records needing defaults
	if count > 0 {
		defaultFields := map[string]interface{}{
			"description": "'No description available'",
		}
		if err := dialect.UpdateNullFields(db, "units", defaultFields, "description IS NULL"); err != nil {
			return fmt.Errorf("failed to set default values for units: %w", err)
		}
	}

	// Only update if there are records needing defaults
	if count > 0 {
		defaultFields := map[string]interface{}{
			"organization_id": 1,
			"branch_id":       "'default'",
		}
		if err := dialect.UpdateNullFields(db, "units", defaultFields, "organization_id IS NULL OR branch_id IS NULL"); err != nil {
			return fmt.Errorf("failed to set default values: %w", err)
		}
	}
	return nil
}

// Rollback: Drop units table
func rollbackUnitsTable(db *gorm.DB, dialect DatabaseDialect) error {
	return db.Migrator().DropTable(&domain.Unit{})
}
