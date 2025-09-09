-- Drop the Triger first
DROP TRIGGER IF EXISTS update_customers_updated_at ON customers;

-- Drop the index
DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_customers_tenant_id;
DROP INDEX IF EXISTS idx_customers_external_id;

-- Drop the table
DROP TABLE IF EXISTS customers;