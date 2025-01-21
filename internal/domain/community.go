package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidCommunityName           = errors.New("community name is required")
	ErrInvalidCommunityAddress        = errors.New("community address is required")
	ErrInvalidCommunityOrganizationID = errors.New("organization ID is required")
	ErrInvalidCommunityBranchID       = errors.New("branch ID is required")
)

type Community struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid"`         // UUID primary key
	Name           string         `json:"name" gorm:"not null"`                   // Community name (not null)
	Address        string         `json:"address" gorm:"not null"`                // Community address (not null)
	Description    string         `json:"description,omitempty" gorm:"type:text"` // Optional description
	CreatedAt      time.Time      `json:"created_at"`                             // Creation timestamp
	UpdatedAt      time.Time      `json:"updated_at"`                             // Update timestamp
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`                         // Soft delete timestamp
	OrganizationID uint           `json:"organization_id" gorm:"not null"`        // Organization ID (required)
	BranchID       string         `json:"branch_id" gorm:"not null"`              // Branch ID (string, required)
}

// Hook to set default values before creating a record
func (c *Community) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Validates the Community object for required fields
func (c *Community) Validate() error {
	if c.Name == "" {
		return ErrInvalidCommunityName
	}
	if c.Address == "" {
		return ErrInvalidCommunityAddress
	}
	if c.OrganizationID == 0 {
		return ErrInvalidCommunityOrganizationID
	}
	if c.BranchID == "" {
		return ErrInvalidCommunityBranchID
	}
	return nil
}

// Request payload for creating a community
type CreateCommunityRequest struct {
	Name           string `json:"name" binding:"required"`
	Address        string `json:"address" binding:"required"`
	Description    string `json:"description,omitempty"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	BranchID       string `json:"branch_id" binding:"required"`
}

// Request payload for updating a community
type UpdateCommunityRequest struct {
	Name           string `json:"name,omitempty"`
	Address        string `json:"address,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID uint   `json:"organization_id,omitempty"`
	BranchID       string `json:"branch_id,omitempty"`
}

// Response payload for a single community
type CommunityResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Description    string    `json:"description,omitempty"`
	OrganizationID uint      `json:"organization_id"`
	BranchID       string    `json:"branch_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Response payload for a list of communities
type ListCommunitiesResponse struct {
	Communities []CommunityResponse `json:"communities"`
	Total       int64               `json:"total"`
}
