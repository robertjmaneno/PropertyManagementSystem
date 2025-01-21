package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/projectname/internal/repository"
	"github.com/yourusername/projectname/pkg/pagination"
	"github.com/yourusername/projectname/pkg/sorting"
	"github.com/yourusername/projectname/test/testutil"
)

func TestRepositoryIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup
	db := testutil.NewTestDB(t)
	userRepo := repository.NewUserRepository(db)

	t.Run("Full User Lifecycle", func(t *testing.T) {
		ctx := context.Background()

		// Create multiple users
		users := testutil.NewTestUsers(10)
		for _, user := range users {
			err := userRepo.Create(ctx, user)
			require.NoError(t, err)
		}

		// Test pagination
		opts := repository.ListOptions{
			Pagination: &pagination.Params{
				Page:     1,
				PageSize: 5,
			},
			Sort: sorting.NewOptions("email", sorting.ASC),
		}

		// Get first page
		result, err := userRepo.GetPaged(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 5)

		// Get total count
		total, err := userRepo.Count(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(10), total)

		// Test search functionality
		searchOpts := repository.ListOptions{
			Pagination: &pagination.Params{
				Page:     1,
				PageSize: 10,
			},
		}

		// Search by email
		searchResults, err := userRepo.Search(ctx, users[0].Email, searchOpts)
		require.NoError(t, err)
		assert.Len(t, searchResults, 1)
		assert.Equal(t, users[0].ID, searchResults[0].ID)

		// Test sorting
		opts.Sort = sorting.NewOptions("email", sorting.DESC)
		sortedResults, err := userRepo.GetPaged(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, sortedResults, 5)
		for i := 1; i < len(sortedResults); i++ {
			assert.True(t, sortedResults[i-1].Email > sortedResults[i].Email)
		}

		// Test update
		user := users[0]
		user.Email = "updated@example.com"
		err = userRepo.Update(ctx, user)
		require.NoError(t, err)

		// Verify update
		updated, err := userRepo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", updated.Email)

		// Test delete
		err = userRepo.Delete(ctx, user.ID)
		require.NoError(t, err)

		// Verify deletion
		total, err = userRepo.Count(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(9), total)
	})
}
