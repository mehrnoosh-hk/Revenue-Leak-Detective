-- Drop trigger first
DROP TRIGGER IF EXISTS update_providers_updated_at ON providers;

-- Drop indexes
DROP INDEX IF EXISTS idx_providers_name;

-- Drop table
DROP TABLE IF EXISTS providers;