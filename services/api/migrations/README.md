Database Migrations
===================
This directory contains the database migration files for the project. Each migration file is named with a sequential three-digit number followed by a descriptive name, indicating the purpose of the migration.

The migrations are managed using golang-migrate, which allows for easy versioning and application of database schema changes.
Migration Files have the following format:
- Up Migration: `NNN_description.up.sql` - Contains SQL statements to apply the migration.
- Down Migration: `NNN_description.down.sql` - Contains SQL statements to revert the migration.

Current migrations:
- 001: Create extensions and functions
- 002: Create tenants table
- 003: Create users table
- 004: Create customers table
- 005: Create leaks table
- 006: Create actions table
- 007: Create payments table
- 008: Create providers table
- 009: Create events table
- 010: Create integrations table
- 011: Add payment_id column to leaks table
The migrations are managed using golang-migrate, which allows for easy versioning and application of database schema changes.
