package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidUnitNumber     = errors.New("unit number is required")
	ErrInvalidUnitRooms      = errors.New("number of rooms is required and must be greater than 0")
	ErrInvalidUnitBathrooms  = errors.New("number of bathrooms is required and must be greater than 0")
	ErrInvalidUnitSquareFeet = errors.New("square feet is required and must be greater than 0")
	ErrInvalidUnitOrgID      = errors.New("organization ID is required")
	ErrInvalidUnitBranchID   = errors.New("branch ID is required")
)

type Unit struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid"`
	UnitNumber     string         `json:"unit_number" gorm:"not null;size:50"`
	Rooms          int            `json:"rooms" gorm:"not null"`
	Bathrooms      float64        `json:"bathrooms" gorm:"not null;type:decimal(1,3)"`
	Description    string         `json:"description,omitempty" gorm:"type:text"`
	SquareFeet     int            `json:"square_feet" gorm:"not null"`
	OrganizationID uint           `json:"organization_id" gorm:"not null"` // Changed to uint
	BranchID       string         `json:"branch_id" gorm:"type:uuid;not null"`
	CreatedAt      time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Unit) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

func (u *Unit) Validate() error {
	if u.UnitNumber == "" {
		return ErrInvalidUnitNumber
	}
	if u.Rooms <= 0 {
		return ErrInvalidUnitRooms
	}
	if u.Bathrooms <= 0 {
		return ErrInvalidUnitBathrooms
	}
	if u.SquareFeet <= 0 {
		return ErrInvalidUnitSquareFeet
	}
	if u.OrganizationID == 0 {
		return ErrInvalidUnitOrgID
	}
	if u.BranchID == "" {
		return ErrInvalidUnitBranchID
	}
	return nil
}

type CreateUnitRequest struct {
	UnitNumber     string  `json:"unit_number" binding:"required"`
	Rooms          int     `json:"rooms" binding:"required"`
	Bathrooms      float64 `json:"bathrooms" binding:"required"`
	Description    string  `json:"description,omitempty"`
	SquareFeet     int     `json:"square_feet" binding:"required"`
	OrganizationID string  `json:"organization_id" binding:"required"`
	BranchID       string  `json:"branch_id" binding:"required"`
}

type UpdateUnitRequest struct {
	UnitNumber     string  `json:"unit_number,omitempty"`
	Rooms          int     `json:"rooms,omitempty"`
	Bathrooms      float64 `json:"bathrooms,omitempty"`
	Description    string  `json:"description,omitempty"`
	SquareFeet     int     `json:"square_feet,omitempty"`
	OrganizationID uint    `json:"organization_id,omitempty"`
	BranchID       string  `json:"branch_id,omitempty"`
}

type UnitResponse struct {
	ID             string    `json:"id"`
	UnitNumber     string    `json:"unit_number"`
	Rooms          int       `json:"rooms"`
	Bathrooms      float64   `json:"bathrooms"`
	Description    string    `json:"description,omitempty"`
	SquareFeet     int       `json:"square_feet"`
	OrganizationID uint      `json:"organization_id"`
	BranchID       string    `json:"branch_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ListUnitsResponse struct {
	Units []UnitResponse `json:"units"`
	Total int64          `json:"total"`
}
