package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yourusername/projectname/internal/config"
	"github.com/yourusername/projectname/internal/migrations"
)

func main() {
	// Parse command line arguments
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Usage: migrate [up|down]")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database using the config method
	db, err := cfg.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	switch args[0] {
	case "up":
		if err := migrations.Migrate(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "down":
		if err := migrations.Rollback(db); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("Rollback completed successfully")
	default:
		log.Fatal("Invalid command. Use 'up' or 'down'")
	}
}
