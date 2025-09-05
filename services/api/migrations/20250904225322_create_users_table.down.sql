-- Drop trigger first
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop index
DROP INDEX IF EXISTS idx_users_email;

-- Drop table
DROP TABLE IF EXISTS users;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";

-- Drop CITEXT extension
DROP EXTENSION IF EXISTS "citext";