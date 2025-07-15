package postgres

import (
	"context"
	"database/sql"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/jmoiron/sqlx"
	"time"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, email, password, role, super_admin) 
        VALUES (:username, :email, :password, :role, :super_admin) 
        RETURNING id
    `
	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&user.ID)
	}
	return sql.ErrNoRows
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, role, super_admin, created_at FROM users WHERE id = $1` // Exclude sensitive fields
	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role, super_admin, created_at FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, user, query, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, role, super_admin, created_at FROM users WHERE username = $1`
	err := r.db.GetContext(ctx, user, query, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users 
        SET username = :username, email = :email, password = :password, role = :role, super_admin = :super_admin 
        WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) ListUsers(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	var users []*models.User
	query := `
        SELECT id, username, email, role, super_admin, created_at 
        FROM users 
        ORDER BY id DESC 
        LIMIT $1 OFFSET $2
    `
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users`
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) UpdateUserRole(ctx context.Context, id int64, role string, isAdmin bool) error {
	query := `UPDATE users SET role = $1, super_admin = $2 WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, role, isAdmin, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, password, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) StoreResetToken(ctx context.Context, id int64, token string) error {
	// Delete existing token if any
	_, _ = r.db.ExecContext(ctx, "DELETE FROM password_reset_tokens WHERE user_id = $1", id)

	query := `
        INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at) 
        VALUES ($1, $2, $3, $4)
    `
	expiresAt := time.Now().Add(time.Hour)
	createdAt := time.Now()
	_, err := r.db.ExecContext(ctx, query, id, token, expiresAt, createdAt)
	return err
}

func (r *userRepository) InvalidateResetToken(ctx context.Context, id int64) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
