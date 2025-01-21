package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidBuildingName     = errors.New("building name is required")
	ErrInvalidBuildingAddress  = errors.New("building address is required")
	ErrInvalidBuildingOrgID    = errors.New("organization ID is required")
	ErrInvalidBuildingBranchID = errors.New("branch ID is required")
)

type Building struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid"`
	Name           string         `json:"name" gorm:"not null"`
	Address        string         `json:"address" gorm:"not null;type:text"`
	Description    string         `json:"description,omitempty" gorm:"type:text"`
	OrganizationID uint           `json:"organization_id" gorm:"not null"` // Changed to uint
	BranchID       string         `json:"branch_id" gorm:"type:uuid;not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (b *Building) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

func (b *Building) Validate() error {
	if b.Name == "" {
		return ErrInvalidBuildingName
	}
	if b.Address == "" {
		return ErrInvalidBuildingAddress
	}
	if b.OrganizationID == 0 { // Check for uint value
		return ErrInvalidBuildingOrgID
	}
	if b.BranchID == "" {
		return ErrInvalidBuildingBranchID
	}
	return nil
}

type CreateBuildingRequest struct {
	Name           string `json:"name" binding:"required"`
	Address        string `json:"address" binding:"required"`
	Description    string `json:"description,omitempty"`
	OrganizationID uint   `json:"organization_id" binding:"required"` // Changed to uint
	BranchID       string `json:"branch_id" binding:"required"`
}

type UpdateBuildingRequest struct {
	Name           string `json:"name,omitempty"`
	Address        string `json:"address,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID uint   `json:"organization_id,omitempty"` // Changed to uint
	BranchID       string `json:"branch_id,omitempty"`
}

type BuildingResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Description    string    `json:"description,omitempty"`
	OrganizationID uint      `json:"organization_id"` // Changed to uint
	BranchID       string    `json:"branch_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ListBuildingsResponse struct {
	Buildings []BuildingResponse `json:"buildings"`
	Total     int64              `json:"total"`
}
