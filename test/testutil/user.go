package testutil

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/projectname/internal/domain"
)

var userCounter = 0

// NewTestUser creates a new user for testing
func NewTestUser() *domain.User {
	userCounter++
	now := time.Now()

	return &domain.User{
		ID:             uuid.New().String(),
		OrganizationID: 1,
		BranchID:       fmt.Sprintf("branch-%d", userCounter),
		Email:          fmt.Sprintf("user%d@example.com", userCounter),
		FirstName:      fmt.Sprintf("First%d", userCounter),
		LastName:       fmt.Sprintf("Last%d", userCounter),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// NewTestUsers creates multiple test users
func NewTestUsers(count int) []*domain.User {
	users := make([]*domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = NewTestUser()
	}
	return users
}
