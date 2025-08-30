-- Migration to make email case-insensitive and ensure uniqueness
-- Migration: 002_email_case_insensitive.up.sql

-- Enable CITEXT extension for case-insensitive text
CREATE EXTENSION IF NOT EXISTS "citext";

-- Drop the existing unique constraint and index
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
DROP INDEX IF EXISTS idx_users_email;

-- Alter the email column to CITEXT type
ALTER TABLE users ALTER COLUMN email TYPE CITEXT;

-- Add a unique constraint on the email column (CITEXT is case-insensitive)
ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);

-- Recreate the index for email lookups (CITEXT indexes are case-insensitive)
CREATE INDEX idx_users_email ON users(email);
