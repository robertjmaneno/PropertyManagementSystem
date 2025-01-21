package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/projectname/internal/domain"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/pkg/pagination"
	"github.com/yourusername/projectname/test/testutil"
)

func TestGenericRepository(t *testing.T) {
	// Setup
	db := testutil.NewTestDB(t)
	repo := repository.NewGenericRepository(db, domain.User{})

	t.Run("Create", func(t *testing.T) {
		user := testutil.NewTestUser()
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Verify
		var found domain.User
		err = db.First(&found, "id = ?", user.ID).Error
		require.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)
	})

	t.Run("GetByID", func(t *testing.T) {
		// Setup
		user := testutil.NewTestUser()
		db.Create(user)

		// Test
		found, err := repo.GetByID(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)

		// Test not found
		_, err = repo.GetByID(context.Background(), "non-existent")
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("GetAll", func(t *testing.T) {
		// Setup
		db.Exec("DELETE FROM users") // Clean up
		users := []*domain.User{
			testutil.NewTestUser(),
			testutil.NewTestUser(),
		}
		for _, u := range users {
			db.Create(u)
		}

		// Test
		found, err := repo.GetAll(context.Background())
		require.NoError(t, err)
		assert.Len(t, found, len(users))
	})

	t.Run("GetPaged", func(t *testing.T) {
		// Setup
		db.Exec("DELETE FROM users") // Clean up
		for i := 0; i < 15; i++ {
			db.Create(testutil.NewTestUser())
		}

		// Test
		opts := repository.ListOptions{
			Pagination: &pagination.Params{
				Page:     1,
				PageSize: 10,
			},
		}

		found, err := repo.GetPaged(context.Background(), opts)
		require.NoError(t, err)
		assert.Len(t, found, 10)

		// Test second page
		opts.Pagination.Page = 2
		found, err = repo.GetPaged(context.Background(), opts)
		require.NoError(t, err)
		assert.Len(t, found, 5)
	})

	t.Run("Update", func(t *testing.T) {
		// Setup
		user := testutil.NewTestUser()
		db.Create(user)

		// Test
		user.Email = "updated@example.com"
		err := repo.Update(context.Background(), user)
		require.NoError(t, err)

		// Verify
		var found domain.User
		err = db.First(&found, "id = ?", user.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", found.Email)
	})

	t.Run("Delete", func(t *testing.T) {
		// Setup
		user := testutil.NewTestUser()
		db.Create(user)

		// Test
		err := repo.Delete(context.Background(), user.ID)
		require.NoError(t, err)

		// Verify
		var found domain.User
		err = db.First(&found, "id = ?", user.ID).Error
		assert.Error(t, err)
	})

	t.Run("Count", func(t *testing.T) {
		// Setup
		db.Exec("DELETE FROM users") // Clean up
		for i := 0; i < 5; i++ {
			db.Create(testutil.NewTestUser())
		}

		// Test
		count, err := repo.Count(context.Background(), repository.ListOptions{})
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}
