package domain_test

import (
	"testing"
	"time"

	"github.com/yourusername/projectname/internal/domain"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    domain.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: domain.User{
				ID:        "user-1",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: domain.User{
				ID:        "user-2",
				Email:     "invalid-email",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: true,
		},
		{
			name: "empty first name",
			user: domain.User{
				ID:        "user-3",
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "User",
			},
			wantErr: true,
		},
		{
			name: "empty last name",
			user: domain.User{
				ID:        "user-4",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.user.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
