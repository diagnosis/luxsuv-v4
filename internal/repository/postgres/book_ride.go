package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/jmoiron/sqlx"
)

type bookRideRepository struct {
	db *sqlx.DB
}

func NewBookRideRepository(db *sqlx.DB) repository.BookRideRepository {
	return &bookRideRepository{db: db}
}

func (r *bookRideRepository) Create(ctx context.Context, br *models.BookRide) error {
	query := `
        INSERT INTO book_rides (user_id, driver_id, your_name, email, phone_number, ride_type, pickup_location, dropoff_location, 
                                date, time, number_of_passengers, number_of_luggage, additional_notes, book_status, ride_status)
        VALUES (:user_id, :driver_id, :your_name, :email, :phone_number, :ride_type, :pickup_location, :dropoff_location, 
                :date, :time, :number_of_passengers, :number_of_luggage, :additional_notes, :book_status, :ride_status)
        RETURNING id
    `
	rows, err := r.db.NamedQueryContext(ctx, query, br)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&br.ID)
	}
	return sql.ErrNoRows
}

func (r *bookRideRepository) GetByID(ctx context.Context, id int64) (*models.BookRide, error) {
	br := &models.BookRide{}
	query := `SELECT * FROM book_rides WHERE id = $1`
	err := r.db.GetContext(ctx, br, query, id)
	if err != nil {
		return nil, err
	}
	return br, nil
}

func (r *bookRideRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.BookRide, error) {
	var bookings []*models.BookRide
	query := `SELECT * FROM book_rides WHERE user_id = $1 ORDER BY date DESC, time DESC`
	err := r.db.SelectContext(ctx, &bookings, query, userID)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookRideRepository) GetByEmail(ctx context.Context, email string) ([]*models.BookRide, error) {
	var bookings []*models.BookRide
	query := `SELECT * FROM book_rides WHERE email = $1 ORDER BY date DESC, time DESC`
	err := r.db.SelectContext(ctx, &bookings, query, email)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookRideRepository) Accept(ctx context.Context, id int64, driverID int64) error {
	query := `
        UPDATE book_rides 
        SET driver_id = $1, book_status = 'Accepted', ride_status = 'Assigned' 
        WHERE id = $2 AND book_status = 'Pending' AND driver_id IS NULL
    `
	result, err := r.db.ExecContext(ctx, query, driverID, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("booking not found, already assigned, or not pending")
	}
	return nil
}
