package uuid

import (
	"github.com/google/uuid"
)

// New generates a new UUID string
func New() string {
	return uuid.New().String()
} 