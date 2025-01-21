package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/yourusername/projectname/docs" // swagger docs
	"github.com/yourusername/projectname/internal/client"
	"github.com/yourusername/projectname/internal/config"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/migrations"
	"github.com/yourusername/projectname/internal/router"
	"github.com/yourusername/projectname/pkg/logger"
)

// Package main is the entry point for the API server.
// It initializes all components and starts the HTTP server.

// @title Template API
// @version 1.0
// @description This is a template API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

// @Security BearerAuth

// AuthClient defines the interface for authentication operations.
// It provides methods for token validation and user login.
type AuthClient interface {
	ValidateToken(ctx context.Context, token string) (*domain.AuthenticatedUser, error)
	Login(ctx context.Context, request domain.LoginRequest) (*domain.AuthResponse, error)
}

func main() {
	// Initialize config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	l := logger.NewLogger(cfg.Log.Level)
	l.Infow("Starting application", "version", "1.0.0")

	// Initialize database
	db, err := cfg.NewDatabase()
	if err != nil {
		l.Fatalw("Failed to connect to database", "error", err)
	}
	l.Info("Connected to database successfully")

	// Run migrations
	if err := migrations.Migrate(db); err != nil {
		l.Fatalw("Failed to migrate database", "error", err)
	}
	l.Info("Database migration completed successfully")

	// Initialize auth client
	authClient := client.NewAuthClient(cfg.Auth.ServiceURL)
	l.Infow("Initialized auth client", "url", cfg.Auth.ServiceURL)

	// Initialize router with dependencies
	r := router.NewRouter(cfg, l, authClient, db)
	l.Info("Router initialized successfully")

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	l.Infow("Starting server", "address", addr)
	if err := r.Run(addr); err != nil {
		l.Fatalw("Failed to start server", "error", err, "address", addr)
	}
}
