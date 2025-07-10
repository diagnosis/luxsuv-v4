package data

import (
	"database/sql"
	"errors"
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
	query := `
INSERT INTO users (username, password, email, role, super_admin, created_at)
VALUES (:username, :password, :email, :role, :super_admin, :created_at)
RETURNING id`
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateUser updates an existing user's fields in the users table.
func (r *Repository) UpdateUser(user *User) error {
	query := `
UPDATE users
SET username = :username,
    password = :password,
    email = :email,
    role = :role,
    super_admin = :super_admin
WHERE id = :id
RETURNING id`
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		return nil
	}
	return sql.ErrNoRows
}

// DeleteUser deletes a user by ID, but only if the requesting user is a super admin.
func (r *Repository) DeleteUser(userID, adminID int64) error {
	// Check if the requesting user is a super admin
	var isSuperAdmin bool
	err := r.db.Get(&isSuperAdmin, "SELECT super_admin FROM users WHERE id = $1", adminID)
	if err != nil {
		return err
	}
	if !isSuperAdmin {
		return errors.New("only super admins can delete users")
	}

	// Delete the user
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *Repository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *Repository) GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
