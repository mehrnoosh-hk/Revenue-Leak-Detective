-- Drop trigger first
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;

-- Drop indexes
DROP INDEX IF EXISTS idx_payments_tenant_id;
DROP INDEX IF EXISTS idx_payments_customer_id;
DROP INDEX IF EXISTS idx_payments_external_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_created_at;

-- Drop table
DROP TABLE IF EXISTS payments;