package e2e

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/handler"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestApp struct {
	DB     *gorm.DB
	Server *httptest.Server
}

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database configuration
	dsn := "host=localhost user=postgres password=postgres dbname=myapp_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Clean up database
	db.Exec("DROP TABLE IF EXISTS users")

	// Run migrations
	err = db.AutoMigrate(&domain.User{})
	require.NoError(t, err)

	return db
}

func setupTestApp(t *testing.T) *TestApp {
	db := setupTestDB(t)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup dependencies
	uow := repository.NewUnitOfWork(db)
	userService := service.NewUserService(uow.Users())
	userHandler := handler.NewUserHandler(userService)

	// Register routes
	userHandler.RegisterRoutes(router.Group("/api"))

	// Create test server
	server := httptest.NewServer(router)

	return &TestApp{
		DB:     db,
		Server: server,
	}
}

func createTestUser(t *testing.T, db *gorm.DB) *domain.User {
	user := &domain.User{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	err := db.Create(user).Error
	require.NoError(t, err)

	return user
}
