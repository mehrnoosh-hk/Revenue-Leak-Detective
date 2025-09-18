-- Drop the Triger first
DROP TRIGGER IF EXISTS update_tenants_updated_at ON tenants;

-- Drop the index
DROP INDEX IF EXISTS idx_tenants_name;

-- Drop the table
DROP TABLE IF EXISTS tenants;