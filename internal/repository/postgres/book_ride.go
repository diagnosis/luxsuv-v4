package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/diagnosis/luxsuv-v4/internal/models"
	"github.com/diagnosis/luxsuv-v4/internal/repository"
	"github.com/jmoiron/sqlx"
	"strings"
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
                                date, time, number_of_passengers, number_of_luggage, additional_notes, book_status, ride_status, created_at, updated_at)
        VALUES (:user_id, :driver_id, :your_name, :email, :phone_number, :ride_type, :pickup_location, :dropoff_location, 
                :date, :time, :number_of_passengers, :number_of_luggage, :additional_notes, :book_status, :ride_status, NOW(), NOW())
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
        SET driver_id = $1, book_status = 'Accepted', ride_status = 'Assigned', updated_at = NOW()
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

func (r *bookRideRepository) Update(ctx context.Context, id int64, updates *models.UpdateBookRideRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.YourName != "" {
		setParts = append(setParts, fmt.Sprintf("your_name = $%d", argIndex))
		args = append(args, updates.YourName)
		argIndex++
	}

	if updates.PhoneNumber != "" {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", argIndex))
		args = append(args, updates.PhoneNumber)
		argIndex++
	}

	if updates.RideType != "" {
		setParts = append(setParts, fmt.Sprintf("ride_type = $%d", argIndex))
		args = append(args, updates.RideType)
		argIndex++
	}

	if updates.PickupLocation != "" {
		setParts = append(setParts, fmt.Sprintf("pickup_location = $%d", argIndex))
		args = append(args, updates.PickupLocation)
		argIndex++
	}

	if updates.DropoffLocation != "" {
		setParts = append(setParts, fmt.Sprintf("dropoff_location = $%d", argIndex))
		args = append(args, updates.DropoffLocation)
		argIndex++
	}

	if updates.Date != "" {
		setParts = append(setParts, fmt.Sprintf("date = $%d", argIndex))
		args = append(args, updates.Date)
		argIndex++
	}

	if updates.Time != "" {
		setParts = append(setParts, fmt.Sprintf("time = $%d", argIndex))
		args = append(args, updates.Time)
		argIndex++
	}

	if updates.NumberOfPassengers != nil {
		setParts = append(setParts, fmt.Sprintf("number_of_passengers = $%d", argIndex))
		args = append(args, *updates.NumberOfPassengers)
		argIndex++
	}

	if updates.NumberOfLuggage != nil {
		setParts = append(setParts, fmt.Sprintf("number_of_luggage = $%d", argIndex))
		args = append(args, *updates.NumberOfLuggage)
		argIndex++
	}

	// Always update additional_notes (can be empty string to clear)
	setParts = append(setParts, fmt.Sprintf("additional_notes = $%d", argIndex))
	args = append(args, updates.AdditionalNotes)
	argIndex++

	if len(setParts) == 0 {
		return errors.New("no fields to update")
	}

	// Always update the updated_at timestamp
	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))

	query := fmt.Sprintf(`
        UPDATE book_rides 
        SET %s
        WHERE id = $%d AND book_status NOT IN ('Cancelled', 'Completed')
    `, strings.Join(setParts, ", "), argIndex)

	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("booking not found or cannot be updated (may be cancelled or completed)")
	}

	return nil
}

func (r *bookRideRepository) Cancel(ctx context.Context, id int64, reason string) error {
	query := `
        UPDATE book_rides 
        SET book_status = 'Cancelled', ride_status = 'Cancelled', additional_notes = COALESCE(additional_notes, '') || $1, updated_at = NOW()
        WHERE id = $2 AND book_status NOT IN ('Cancelled', 'Completed')
    `
	
	cancelNote := fmt.Sprintf("\n[CANCELLED: %s]", reason)
	result, err := r.db.ExecContext(ctx, query, cancelNote, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("booking not found or cannot be cancelled (may already be cancelled or completed)")
	}

	return nil
}

func (r *bookRideRepository) GetByIDAndEmail(ctx context.Context, id int64, email string) (*models.BookRide, error) {
	br := &models.BookRide{}
	query := `SELECT * FROM book_rides WHERE id = $1 AND email = $2`
	err := r.db.GetContext(ctx, br, query, id, email)
	if err != nil {
		return nil, err
	}
	return br, nil
}