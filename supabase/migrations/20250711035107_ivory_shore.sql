-- +goose Up
-- +goose StatementBegin

-- Add unique constraint to username column
ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove indexes
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;

-- Remove unique constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_unique;

-- +goose StatementEnd