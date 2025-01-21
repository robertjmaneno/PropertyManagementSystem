package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/pkg/errors"
)

func TestUserFlow(t *testing.T) {
	app := setupTestApp(t)
	client := NewTestClient(app.Server.URL)

	// Test user creation
	t.Run("create and get user", func(t *testing.T) {
		newUser := &domain.User{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}

		// Create user
		createdUser, err := client.CreateUser(newUser)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		assert.NotEmpty(t, createdUser.ID)
		assert.Equal(t, newUser.Email, createdUser.Email)
		assert.Equal(t, newUser.FirstName, createdUser.FirstName)
		assert.Equal(t, newUser.LastName, createdUser.LastName)

		// Verify user exists in database
		var dbUser domain.User
		err = app.DB.First(&dbUser, "id = ?", createdUser.ID).Error
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, dbUser.ID)
		assert.Equal(t, newUser.Email, dbUser.Email)
		assert.Equal(t, newUser.FirstName, dbUser.FirstName)
		assert.Equal(t, newUser.LastName, dbUser.LastName)

		// Get user via API
		fetchedUser, err := client.GetUser(createdUser.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedUser)
		assert.Equal(t, createdUser.ID, fetchedUser.ID)
		assert.Equal(t, newUser.Email, fetchedUser.Email)
		assert.Equal(t, newUser.FirstName, fetchedUser.FirstName)
		assert.Equal(t, newUser.LastName, fetchedUser.LastName)
	})

	// Test user not found
	t.Run("get non-existent user", func(t *testing.T) {
		user, err := client.GetUser("non-existent")
		assert.Error(t, err)
		assert.Equal(t, errors.ErrNotFound, err)
		assert.Nil(t, user)
	})

	// Test invalid user creation
	t.Run("create invalid user", func(t *testing.T) {
		invalidUser := &domain.User{
			Email:     "invalid-email",
			FirstName: "",
			LastName:  "",
		}

		user, err := client.CreateUser(invalidUser)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
