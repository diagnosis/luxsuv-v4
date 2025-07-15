package repository

import (
	"context"
	"github.com/diagnosis/luxsuv-v4/internal/models"
)

type BookRideRepository interface {
	Create(ctx context.Context, br *models.BookRide) error
	GetByID(ctx context.Context, id int64) (*models.BookRide, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.BookRide, error)
	GetByEmail(ctx context.Context, email string) ([]*models.BookRide, error)
	Accept(ctx context.Context, id int64, driverID int64) error
}
