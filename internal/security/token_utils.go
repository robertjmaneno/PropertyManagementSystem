package security

import (
	"strings"
)

const (
	Bearer = "Bearer "
)

// ExtractBearerToken extracts the bearer token from the Authorization header
func ExtractBearerToken(authHeader string) string {
	if authHeader != "" && strings.HasPrefix(authHeader, Bearer) {
		return strings.TrimPrefix(authHeader, Bearer)
	}
	return ""
}

// CreateBearerToken creates a bearer token string from a token
func CreateBearerToken(token string) string {
	return Bearer + token
} 