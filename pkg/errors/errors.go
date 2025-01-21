package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code       string                 // Error code for client handling
	Message    string                 // User-friendly error message
	Details    map[string]interface{} // Additional error details
	Err        error                  // Original error
	StatusCode int                    // HTTP status code
}

func (e *AppError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("%s: %s (details: %v)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code string, message string, err error, statusCode int, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		Err:        err,
		StatusCode: statusCode,
	}
}

// Common application errors
var (
	// Generic errors
	ErrNotFound = NewAppError(
		"NOT_FOUND",
		"Resource not found",
		nil,
		http.StatusNotFound,
		nil,
	)
	ErrInvalidInput = NewAppError(
		"INVALID_INPUT",
		"Invalid input provided",
		nil,
		http.StatusBadRequest,
		nil,
	)
	ErrInternalServer = NewAppError(
		"INTERNAL_SERVER_ERROR",
		"Internal server error",
		nil,
		http.StatusInternalServerError,
		nil,
	)
	ErrDuplicateResource = NewAppError(
		"DUPLICATE_RESOURCE",
		"Resource already exists",
		nil,
		http.StatusConflict,
		nil,
	)

	// Authentication errors
	ErrInvalidToken = NewAppError(
		"INVALID_TOKEN",
		"The provided authentication token is invalid",
		nil,
		http.StatusUnauthorized,
		nil,
	)
	ErrExpiredToken = NewAppError(
		"EXPIRED_TOKEN",
		"The authentication token has expired",
		nil,
		http.StatusUnauthorized,
		nil,
	)
	ErrMissingToken = NewAppError(
		"MISSING_TOKEN",
		"No authentication token was provided",
		nil,
		http.StatusUnauthorized,
		nil,
	)
	ErrInvalidCredentials = NewAppError(
		"INVALID_CREDENTIALS",
		"The provided credentials are invalid",
		nil,
		http.StatusUnauthorized,
		nil,
	)

	// Authorization errors
	ErrUnauthorized = NewAppError(
		"UNAUTHORIZED",
		"User is not authorized to perform this action",
		nil,
		http.StatusUnauthorized,
		nil,
	)
	ErrInsufficientPermissions = NewAppError(
		"INSUFFICIENT_PERMISSIONS",
		"User does not have sufficient permissions",
		nil,
		http.StatusForbidden,
		nil,
	)
	ErrAccountDisabled = NewAppError(
		"ACCOUNT_DISABLED",
		"The user account is disabled",
		nil,
		http.StatusForbidden,
		nil,
	)
	ErrAccountLocked = NewAppError(
		"ACCOUNT_LOCKED",
		"The user account is locked",
		nil,
		http.StatusForbidden,
		nil,
	)

	// CSRF errors
	ErrInvalidCSRFToken = NewAppError(
		"INVALID_CSRF_TOKEN",
		"Invalid or missing CSRF token",
		nil,
		http.StatusForbidden,
		nil,
	)
	ErrCSRFTokenExpired = NewAppError(
		"CSRF_TOKEN_EXPIRED",
		"The CSRF token has expired",
		nil,
		http.StatusForbidden,
		nil,
	)

	// Service errors
	ErrAuthServiceUnavailable = NewAppError(
		"AUTH_SERVICE_UNAVAILABLE",
		"The authentication service is currently unavailable",
		nil,
		http.StatusServiceUnavailable,
		nil,
	)
	ErrPermissionServiceUnavailable = NewAppError(
		"PERMISSION_SERVICE_UNAVAILABLE",
		"The permission service is currently unavailable",
		nil,
		http.StatusServiceUnavailable,
		nil,
	)

	// Session errors
	ErrInvalidSession = NewAppError(
		"INVALID_SESSION",
		"The session is invalid or has expired",
		nil,
		http.StatusUnauthorized,
		nil,
	)
	ErrConcurrentSession = NewAppError(
		"CONCURRENT_SESSION",
		"Another session has been started with this account",
		nil,
		http.StatusConflict,
		nil,
	)

	ErrDuplicateEntry = NewAppError(
		"DUPLICATE_ENTRY",
		"A duplicate entry was found",
		nil,
		http.StatusConflict,
		nil,
	)

	ErrForbidden = errors.New("forbidden")
)

// WrapError wraps an error with additional context
func WrapError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound
	case errors.Is(err, ErrInvalidInput):
		return ErrInvalidInput
	case errors.Is(err, ErrUnauthorized):
		return ErrUnauthorized
	case errors.Is(err, ErrDuplicateResource):
		return ErrDuplicateResource
	case errors.Is(err, ErrDuplicateEntry):
		return ErrDuplicateEntry
	default:
		return ErrInternalServer
	}
}

// WithDetails adds details to an AppError
func WithDetails(err *AppError, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       err.Code,
		Message:    err.Message,
		Details:    details,
		Err:        err.Err,
		StatusCode: err.StatusCode,
	}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized checks if the error is an unauthorized error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsInvalidInput checks if the error is an invalid input error
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}
