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
	Update(ctx context.Context, id int64, updates *models.UpdateBookRideRequest) error
	Cancel(ctx context.Context, id int64, reason string) error
	GetByIDAndEmail(ctx context.Context, id int64, email string) (*models.BookRide, error)
	GetAvailableBookings(ctx context.Context) ([]*models.BookingListResponse, error)
	GetAssignedBookings(ctx context.Context, driverID int64) ([]*models.BookingListResponse, error)
	GetAllBookingsForDispatcher(ctx context.Context) ([]*models.BookingListResponse, error)
	AssignToDriver(ctx context.Context, bookingID int64, driverID int64, assignedBy int64, notes string) error
	GetDriverBookings(ctx context.Context, driverID int64, status string) ([]*models.BookingListResponse, error)
}
