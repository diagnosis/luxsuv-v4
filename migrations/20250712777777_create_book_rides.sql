-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE book_rides (
                            id BIGSERIAL PRIMARY KEY,
                            user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
                            driver_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
                            your_name TEXT NOT NULL,
                            email TEXT NOT NULL,
                            phone_number TEXT NOT NULL,
                            ride_type TEXT NOT NULL,
                            pickup_location TEXT NOT NULL,
                            dropoff_location TEXT NOT NULL,
                            date TEXT NOT NULL,
                            time TEXT NOT NULL,
                            number_of_passengers INTEGER NOT NULL,
                            number_of_luggage INTEGER NOT NULL,
                            additional_notes TEXT,
                            book_status TEXT NOT NULL DEFAULT 'Pending',
                            ride_status TEXT NOT NULL DEFAULT 'Pending'
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE book_rides;
-- +goose StatementEnd
