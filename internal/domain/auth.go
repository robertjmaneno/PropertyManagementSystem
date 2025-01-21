package domain

// AuthenticatedUser represents a user authenticated by the auth service
type AuthenticatedUser struct {
	EmployeeID   string `json:"employeeId"`
	EmailAddress string `json:"emailAddress"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	PhoneNumber  string `json:"phoneNumber"`
	Enabled      bool   `json:"enabled"`
	Roles        []Role `json:"roles"`
}

// Role represents a user role
type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// LoginRequest represents the login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token        string       `json:"token"`
	AccessTicket AccessTicket `json:"accessTicket"`
	IAT          float64      `json:"iat"`
	EXP          float64      `json:"exp"`
}

// AccessTicket represents the access ticket in the auth response
type AccessTicket struct {
	Sub          string `json:"sub"`
	Username     string `json:"username"`
	EmployeeID   string `json:"employeeId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	PhoneNumber  string `json:"phoneNumber"`
	Enabled      bool   `json:"enabled"`
	Roles        []Role `json:"roles"`
	PendingReset bool   `json:"pendingReset"`
}

// Common errors
var (
	ErrInvalidToken        = NewAuthError("invalid token")
	ErrMissingToken        = NewAuthError("missing token")
	ErrInvalidCredentials  = NewAuthError("invalid credentials")
	ErrUnauthorized        = NewAuthError("unauthorized")
	ErrAuthServiceUnavailable = NewAuthError("auth service unavailable")
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

func NewAuthError(message string) error {
	return &AuthError{Message: message}
} 