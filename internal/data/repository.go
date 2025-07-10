package data

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Repository manages database operations for users.
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance with the provided database connection.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a new user into the users table and returns the generated ID.
func (r *Repository) CreateUser(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	query := `
INSERT INTO users (username, password, email, role, super_admin, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`
	
	err := r.db.QueryRow(query, user.Username, user.Password, user.Email, user.Role, user.SuperAdmin, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// UpdateUser updates an existing user's fields in the users table.
func (r *Repository) UpdateUser(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	
	if user.ID <= 0 {
		return errors.New("invalid user ID")
	}

	query := `
UPDATE users
SET username = $1,
    password = $2,
    email = $3,
    role = $4,
    super_admin = $5
WHERE id = $6`
	
	result, err := r.db.Exec(query, user.Username, user.Password, user.Email, user.Role, user.SuperAdmin, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteUser deletes a user by ID, but only if the requesting user is a super admin.
func (r *Repository) DeleteUser(userID, adminID int64) error {
	if userID <= 0 || adminID <= 0 {
		return errors.New("invalid user ID or admin ID")
	}

	// Check if the requesting user is a super admin
	var isSuperAdmin bool
	err := r.db.Get(&isSuperAdmin, "SELECT super_admin FROM users WHERE id = $1", adminID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("admin user not found")
		}
		return fmt.Errorf("failed to check admin privileges: %w", err)
	}
	
	if !isSuperAdmin {
		return errors.New("only super admins can delete users")
	}

	// Prevent self-deletion
	if userID == adminID {
		return errors.New("cannot delete your own account")
	}

	// Delete the user
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *Repository) GetUserByEmail(email string) (*User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user := &User{}
	query := "SELECT id, username, password, email, role, super_admin, created_at FROM users WHERE email = $1"
	
	err := r.db.Get(user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *Repository) GetUserByID(id int64) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user := &User{}
	query := "SELECT id, username, password, email, role, super_admin, created_at FROM users WHERE id = $1"
	
	err := r.db.Get(user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *Repository) GetUserByUsername(username string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	user := &User{}
	query := "SELECT id, username, password, email, role, super_admin, created_at FROM users WHERE username = $1"
	
	err := r.db.Get(user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	
	return user, nil
}

// ListUsers retrieves all users with pagination support.
func (r *Repository) ListUsers(limit, offset int) ([]*User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	var users []*User
	query := `
SELECT id, username, password, email, role, super_admin, created_at 
FROM users 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2`
	
	err := r.db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	
	return users, nil
}

// CountUsers returns the total number of users.
func (r *Repository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users")
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}