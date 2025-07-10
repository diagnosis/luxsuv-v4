package repository

import (
	"database/sql"
	"fmt"

	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, role, super_admin, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at`

	err := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role, user.IsAdmin).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password, role, super_admin, created_at
		FROM users 
		WHERE email = $1`

	err := r.db.Get(user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password, role, super_admin, created_at
		FROM users 
		WHERE id = $1`

	err := r.db.Get(user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password, role, super_admin, created_at
		FROM users 
		WHERE username = $1`

	err := r.db.Get(user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(userID int64) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
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

// ListUsers retrieves all users with pagination
func (r *UserRepository) ListUsers(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := `
		SELECT id, username, email, role, super_admin, created_at
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	err := r.db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users")
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// UpdateUserRole updates a user's role and admin status
func (r *UserRepository) UpdateUserRole(userID int64, role string, isAdmin bool) error {
	query := `UPDATE users SET role = $1, super_admin = $2 WHERE id = $3`
	
	result, err := r.db.Exec(query, role, isAdmin, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
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

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID int64, hashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	
	result, err := r.db.Exec(query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
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

// StoreResetToken stores a password reset token for a user
func (r *UserRepository) StoreResetToken(userID int64, token string) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, NOW() + INTERVAL '1 hour', NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET token = $2, expires_at = NOW() + INTERVAL '1 hour', created_at = NOW()`
	
	_, err := r.db.Exec(query, userID, token)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	return nil
}

// GetResetToken retrieves a reset token for validation
func (r *UserRepository) GetResetToken(token string) (int64, error) {
	var userID int64
	query := `
		SELECT user_id FROM password_reset_tokens 
		WHERE token = $1 AND expires_at > NOW()`
	
	err := r.db.Get(&userID, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, sql.ErrNoRows
		}
		return 0, fmt.Errorf("failed to get reset token: %w", err)
	}

	return userID, nil
}

// InvalidateResetToken removes a reset token
func (r *UserRepository) InvalidateResetToken(userID int64) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to invalidate reset token: %w", err)
	}

	return nil
}