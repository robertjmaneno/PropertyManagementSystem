package fixtures

import (
	"time"

	"github.com/yourusername/projectname/internal/domain"
)

// TestUser1 is a test user fixture
var TestUser1 = &domain.User{
	ID:             "user-1",
	OrganizationID: 1,
	BranchID:       "branch-1",
	Email:          "test1@example.com",
	FirstName:      "Test",
	LastName:       "User 1",
	CreatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	UpdatedAt:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
}

// TestUser2 is another test user fixture
var TestUser2 = &domain.User{
	ID:             "user-2",
	OrganizationID: 1,
	BranchID:       "branch-1",
	Email:          "test2@example.com",
	FirstName:      "Test",
	LastName:       "User 2",
	CreatedAt:      time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	UpdatedAt:      time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
}

// CreateTestUsers creates test users in the database
func CreateTestUsers(queries []string) []string {
	return []string{
		`INSERT INTO users (id, organization_id, branch_id, email, first_name, last_name, created_at, updated_at)
		VALUES ('user-1', 1, 'branch-1', 'test1@example.com', 'Test', 'User 1',
		'2023-01-01 00:00:00+00', '2023-01-01 00:00:00+00')`,

		`INSERT INTO users (id, organization_id, branch_id, email, first_name, last_name, created_at, updated_at)
		VALUES ('user-2', 1, 'branch-1', 'test2@example.com', 'Test', 'User 2',
		'2023-01-02 00:00:00+00', '2023-01-02 00:00:00+00')`,
	}
}
