-- Drop the Triger first
DROP TRIGGER IF EXISTS update_integrations_updated_at ON integrations;

-- Drop the index
DROP INDEX IF EXISTS idx_integrations_tenant_id;
DROP INDEX IF EXISTS idx_integrations_provider_id;
DROP INDEX IF EXISTS idx_integrations_created_at;
DROP INDEX IF EXISTS idx_integrations_updated_at;

-- Drop the table
DROP TABLE IF EXISTS integrations;