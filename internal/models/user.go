package models

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Never include in JSON responses
	Role      string    `json:"role" db:"role"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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