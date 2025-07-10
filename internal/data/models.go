package data

import "time"

// User represents a user entity (e.g., driver, customer, or super admin) in the LuxSUV system.
type User struct {
	ID         int64     `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"password" db:"password"`
	Email      string    `json:"email" db:"email"`
	Role       string    `json:"role" db:"role"`
	SuperAdmin bool      `json:"super_admin" db:"super_admin"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	Token      string    `json:"token" db:"-"` // Not stored in DB
}
