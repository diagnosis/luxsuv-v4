package models

import (
	"time"
)

type User struct {
	ID         int       `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"-" db:"password"` // Hide password in JSON responses
	Email      string    `json:"email" db:"email"`
	Role       string    `json:"role" db:"role"`
	SuperAdmin bool      `json:"super_admin" db:"super_admin"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}