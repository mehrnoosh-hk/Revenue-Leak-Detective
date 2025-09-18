-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop all ENUM types (in reverse dependency order)
-- Note: These will only be dropped if no tables are using them
DROP TYPE IF EXISTS event_status_enum CASCADE;
DROP TYPE IF EXISTS event_type_enum CASCADE;
DROP TYPE IF EXISTS payment_type_enum CASCADE;
DROP TYPE IF EXISTS payment_status_enum CASCADE;
DROP TYPE IF EXISTS action_result_enum CASCADE;
DROP TYPE IF EXISTS action_status_enum CASCADE;
DROP TYPE IF EXISTS action_type_enum CASCADE;
DROP TYPE IF EXISTS leak_type_enum CASCADE;

-- Drop CITEXT extension
DROP EXTENSION IF EXISTS "citext";

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
