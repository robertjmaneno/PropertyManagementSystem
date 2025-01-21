package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/test/testutil"
)

func TestUserRepository_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	db, err := testutil.NewTestDB()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Auto migrate the test models
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	// Create repository
	repo := repository.NewUserRepository(db)

	// Test cases
	t.Run("Create and retrieve user", func(t *testing.T) {
		ctx := context.Background()

		// Create test user
		user := &domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}

		// Create user
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)

		// Retrieve user
		retrieved, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.FirstName, retrieved.FirstName)
		assert.Equal(t, user.LastName, retrieved.LastName)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	// Clean up
	db.Migrator().DropTable(&domain.User{})
}
