package domain

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Common errors
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrEmptyFirstName        = errors.New("first name cannot be empty")
	ErrEmptyLastName         = errors.New("last name cannot be empty")
	ErrInvalidOrganizationID = errors.New("organization ID is required")
	ErrInvalidBranchID       = errors.New("branch ID is required")
)

// User represents a user in the system
type User struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid"`
	FirstName      string         `json:"first_name" gorm:"not null"`
	LastName       string         `json:"last_name" gorm:"not null"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null"`
	OrganizationID uint           `json:"organization_id" gorm:"not null"`
	BranchID       string         `json:"branch_id" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Validate checks if the user data is valid
func (u *User) Validate() error {
	// Check email
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return ErrInvalidEmail
	}

	// Check first name
	if u.FirstName == "" {
		return ErrEmptyFirstName
	}

	// Check last name
	if u.LastName == "" {
		return ErrEmptyLastName
	}

	// Check organization ID
	if u.OrganizationID == 0 {
		return ErrInvalidOrganizationID
	}

	// Check branch ID
	if u.BranchID == "" {
		return ErrInvalidBranchID
	}

	return nil
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email          string   `json:"email" binding:"required,email"`
	FirstName      string   `json:"firstName" binding:"required"`
	LastName       string   `json:"lastName" binding:"required"`
	OrganizationID uint     `json:"organizationId" binding:"required"`
	BranchID       string   `json:"branchId" binding:"required"`
	Roles          []string `json:"roles" binding:"required"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	FirstName      string   `json:"firstName,omitempty"`
	LastName       string   `json:"lastName,omitempty"`
	OrganizationID uint     `json:"organizationId,omitempty"`
	BranchID       string   `json:"branchId,omitempty"`
	Roles          []string `json:"roles,omitempty"`
}

// UserResponse represents the response body for user endpoints
type UserResponse struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	OrganizationID uint      `json:"organizationId"`
	BranchID       string    `json:"branchId"`
	Roles          []Role    `json:"roles"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// ListUsersResponse represents the response body for listing users
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}
