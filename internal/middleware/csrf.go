package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/projectname/pkg/errors"
)

const (
	// CSRF token length in bytes
	csrfTokenLength = 32
	// CSRF token expiry duration
	csrfTokenExpiry = 24 * time.Hour
	// CSRF token header name
	csrfTokenHeader = "X-CSRF-Token"
	// CSRF token cookie name
	csrfTokenCookie = "csrf_token"
)

// CSRFStore stores CSRF tokens and their expiry times
type CSRFStore struct {
	tokens map[string]time.Time
	mutex  sync.RWMutex
}

// NewCSRFStore creates a new CSRF token store
func NewCSRFStore() *CSRFStore {
	store := &CSRFStore{
		tokens: make(map[string]time.Time),
	}
	
	// Start cleanup goroutine
	go store.cleanup()
	
	return store
}

// cleanup periodically removes expired tokens
func (s *CSRFStore) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for token, expiry := range s.tokens {
			if now.After(expiry) {
				delete(s.tokens, token)
			}
		}
		s.mutex.Unlock()
	}
}

// generateToken generates a new CSRF token
func generateToken() (string, error) {
	bytes := make([]byte, csrfTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CSRFMiddleware creates a middleware for CSRF protection
func CSRFMiddleware(store *CSRFStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for safe methods
		if isSafeMethod(c.Request.Method) {
			// Generate and set new token for GET requests
			if c.Request.Method == http.MethodGet {
				token, err := generateToken()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": errors.ErrInternalServer.Error(),
					})
					return
				}

				// Store token with expiry
				store.mutex.Lock()
				store.tokens[token] = time.Now().Add(csrfTokenExpiry)
				store.mutex.Unlock()

				// Set token in cookie and header
				c.SetCookie(csrfTokenCookie, token, int(csrfTokenExpiry.Seconds()), "/", "", true, true)
				c.Header(csrfTokenHeader, token)
			}
			c.Next()
			return
		}

		// Verify CSRF token for unsafe methods
		token := c.GetHeader(csrfTokenHeader)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.ErrInvalidCSRFToken.Error(),
			})
			return
		}

		// Verify token exists and hasn't expired
		store.mutex.RLock()
		expiry, exists := store.tokens[token]
		
		store.mutex.RUnlock()

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.ErrInvalidCSRFToken.Error(),
			})
			return
		}

		if time.Now().After(expiry) {
			store.mutex.Lock()
			delete(store.tokens, token)
			store.mutex.Unlock()
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": errors.ErrCSRFTokenExpired.Error(),
			})
			return
		}

		c.Next()
	}
}

// isSafeMethod checks if the HTTP method is safe (doesn't modify state)
func isSafeMethod(method string) bool {
	method = strings.ToUpper(method)
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions
} 