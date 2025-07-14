package repository

import (
	"context"
	"github.com/diagnosis/luxsuv-v4/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error) // Note: Standardized to GetByID for consistency
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, limit int, offset int) ([]*models.User, error) // Use pointers for efficiency
	CountUsers(ctx context.Context) (int64, error)
	UpdateUserRole(ctx context.Context, id int64, role string, isAdmin bool) error
	UpdatePassword(ctx context.Context, id int64, password string) error
	StoreResetToken(ctx context.Context, id int64, token string) error
	InvalidateResetToken(ctx context.Context, id int64) error
}
