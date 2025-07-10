package models

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Never include in JSON responses
	Role      string    `json:"role" db:"role"`
	IsAdmin   bool      `json:"is_admin" db:"super_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserRole constants
const (
	RoleRider  = "rider"
	RoleDriver = "driver"
	RoleAdmin  = "admin"
)

// IsValidRole checks if the role is valid
func IsValidRole(role string) bool {
	switch role {
	case RoleRider, RoleDriver, RoleAdmin:
		return true
	default:
		return false
	}
}

// CreateUserRequest represents the request payload for user registration
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=rider driver admin"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}