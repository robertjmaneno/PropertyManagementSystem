package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/projectname/internal/domain"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserResponse represents the response body for user operations
type UserResponse struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	OrganizationID uint      `json:"organization_id"`
	BranchID       string    `json:"branch_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ListUsersResponse represents the response body for listing users
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

// ToUser converts a CreateUserRequest to a domain.User
func (r *CreateUserRequest) ToUser() *domain.User {
	return &domain.User{
		ID:        uuid.New().String(),
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
	}
}

// FromUser creates a UserResponse from a domain.User
func FromUser(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		OrganizationID: user.OrganizationID,
		BranchID:       user.BranchID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

// FromUsers creates a ListUsersResponse from a slice of domain.User
func FromUsers(users []domain.User, total int64) *ListUsersResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		response := FromUser(&user)
		responses[i] = *response
	}

	return &ListUsersResponse{
		Users: responses,
		Total: total,
	}
}
