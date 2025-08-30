Database Migrations
===================
This directory contains the database migration files for the project. Each migration file is named with a number followed by a descriptive name, indicating the purpose of the migration.
Migration Files have the following format:
- Up Migration: `XXXX_description.up.sql` - Contains SQL statements to apply the migration.
- Down Migration: `XXXX_description.down.sql` - Contains SQL statements to revert the migration.
The migrations are managed using golang-migrate, which allows for easy versioning and application of database schema changes.
