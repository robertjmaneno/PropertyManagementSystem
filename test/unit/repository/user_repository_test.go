package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/errors"
	"github.com/yourusername/projectname/pkg/pagination"
	"github.com/yourusername/projectname/pkg/sorting"
	"github.com/yourusername/projectname/test/testutil"
)

func TestUserRepository(t *testing.T) {
	// Setup
	db := testutil.NewTestDB(t)
	repo := repository.NewUserRepository(db)

	t.Run("GetByEmail", func(t *testing.T) {
		// Setup
		user := testutil.NewTestUser()
		db.Create(user)

		// Test success
		found, err := repo.GetByEmail(context.Background(), user.Email)
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, user.Email, found.Email)

		// Test not found
		_, err = repo.GetByEmail(context.Background(), "nonexistent@example.com")
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("Search", func(t *testing.T) {
		// Setup
		db.Exec("DELETE FROM users") // Clean up
		users := testutil.NewTestUsers(5)
		for _, u := range users {
			db.Create(u)
		}

		// Test search by email
		opts := repository.ListOptions{
			Pagination: &pagination.Params{
				Page:     1,
				PageSize: 10,
			},
			Sort: sorting.NewOptions("email", sorting.ASC),
		}

		found, err := repo.Search(context.Background(), users[0].Email, opts)
		require.NoError(t, err)
		assert.Len(t, found, 1)
		assert.Equal(t, users[0].ID, found[0].ID)

		// Test search by first name
		found, err = repo.Search(context.Background(), users[1].FirstName, opts)
		require.NoError(t, err)
		assert.Len(t, found, 1)
		assert.Equal(t, users[1].ID, found[0].ID)

		// Test search by last name
		found, err = repo.Search(context.Background(), users[2].LastName, opts)
		require.NoError(t, err)
		assert.Len(t, found, 1)
		assert.Equal(t, users[2].ID, found[0].ID)

		// Test search with no results
		found, err = repo.Search(context.Background(), "nonexistent", opts)
		require.NoError(t, err)
		assert.Len(t, found, 0)

		// Test search with pagination
		opts.Pagination.PageSize = 2
		found, err = repo.Search(context.Background(), "user", opts)
		require.NoError(t, err)
		assert.Len(t, found, 2)

		// Test search with sorting
		opts.Sort = sorting.NewOptions("email", sorting.DESC)
		found, err = repo.Search(context.Background(), "user", opts)
		require.NoError(t, err)
		assert.Len(t, found, 2)
		assert.True(t, found[0].Email > found[1].Email)
	})

	t.Run("Generic Repository Methods", func(t *testing.T) {
		// Test that inherited methods work
		user := testutil.NewTestUser()

		// Test Create
		err := repo.Create(context.Background(), user)
		require.NoError(t, err)

		// Test GetByID
		found, err := repo.GetByID(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)

		// Test Update
		user.Email = "updated@example.com"
		err = repo.Update(context.Background(), user)
		require.NoError(t, err)

		// Verify update
		found, err = repo.GetByID(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", found.Email)

		// Test Delete
		err = repo.Delete(context.Background(), user.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(context.Background(), user.ID)
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})
}
