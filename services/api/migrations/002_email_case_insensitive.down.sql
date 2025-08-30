-- Rollback migration for email case-insensitive changes
-- Migration: 002_email_case_insensitive.down.sql

-- Drop the unique constraint and index
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
DROP INDEX IF EXISTS idx_users_email;

-- Alter the email column back to VARCHAR(255)
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(255);

-- Add back the original unique constraint (case-sensitive)
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

-- Recreate the original index
CREATE INDEX idx_users_email ON users(email);

-- Note: We don't drop the CITEXT extension as it might be used by other parts of the system
