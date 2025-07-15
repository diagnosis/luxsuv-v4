-- +goose Up
-- +goose StatementBegin

/*
  # Update book_rides table for update/cancel functionality

  1. New Columns
    - `created_at` (timestamp) - Track when booking was created
    - `updated_at` (timestamp) - Track when booking was last updated

  2. Indexes
    - Add indexes for better query performance on email, status, and timestamps

  3. Notes
    - Adds timestamp tracking for audit purposes
    - Improves query performance with strategic indexes
*/

-- Add timestamp columns
ALTER TABLE book_rides 
ADD COLUMN created_at TIMESTAMP DEFAULT NOW(),
ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();

-- Update existing records to have proper timestamps
UPDATE book_rides SET created_at = NOW(), updated_at = NOW() WHERE created_at IS NULL;

-- Make timestamp columns NOT NULL after setting defaults
ALTER TABLE book_rides 
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_book_rides_email ON book_rides(email);
CREATE INDEX IF NOT EXISTS idx_book_rides_user_id ON book_rides(user_id);
CREATE INDEX IF NOT EXISTS idx_book_rides_driver_id ON book_rides(driver_id);
CREATE INDEX IF NOT EXISTS idx_book_rides_book_status ON book_rides(book_status);
CREATE INDEX IF NOT EXISTS idx_book_rides_ride_status ON book_rides(ride_status);
CREATE INDEX IF NOT EXISTS idx_book_rides_date ON book_rides(date);
CREATE INDEX IF NOT EXISTS idx_book_rides_created_at ON book_rides(created_at);
CREATE INDEX IF NOT EXISTS idx_book_rides_updated_at ON book_rides(updated_at);

-- Composite index for common queries
CREATE INDEX IF NOT EXISTS idx_book_rides_email_status ON book_rides(email, book_status);
CREATE INDEX IF NOT EXISTS idx_book_rides_user_status ON book_rides(user_id, book_status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove indexes
DROP INDEX IF EXISTS idx_book_rides_user_status;
DROP INDEX IF EXISTS idx_book_rides_email_status;
DROP INDEX IF EXISTS idx_book_rides_updated_at;
DROP INDEX IF EXISTS idx_book_rides_created_at;
DROP INDEX IF EXISTS idx_book_rides_date;
DROP INDEX IF EXISTS idx_book_rides_ride_status;
DROP INDEX IF EXISTS idx_book_rides_book_status;
DROP INDEX IF EXISTS idx_book_rides_driver_id;
DROP INDEX IF EXISTS idx_book_rides_user_id;
DROP INDEX IF EXISTS idx_book_rides_email;

-- Remove timestamp columns
ALTER TABLE book_rides 
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS created_at;

-- +goose StatementEnd