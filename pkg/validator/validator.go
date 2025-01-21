package validator

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Custom validation tags
const (
	passwordMinLen = 8
)

// Setup initializes the validator with custom validations
func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register custom validation tags
		v.RegisterValidation("password", validatePassword)
		
		// Register custom error messages
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})
	}
}

// validatePassword checks if the password meets the requirements
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Check minimum length
	if len(password) < passwordMinLen {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune(`!@#$%^&*()_+-=[]{}|;:,.<>?`, char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Validate validates a struct using tags
func Validate(i interface{}) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v.Struct(i)
	}
	return nil
} 