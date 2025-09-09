-- Drop trigger first
DROP TRIGGER IF EXISTS update_leaks_updated_at ON leaks;

-- Drop indexes
DROP INDEX IF EXISTS idx_leaks_tenant_id;
DROP INDEX IF EXISTS idx_leaks_customer_id;
DROP INDEX IF EXISTS idx_leaks_payment_id;
DROP INDEX IF EXISTS idx_leaks_type;
DROP INDEX IF EXISTS idx_leaks_created_at;

-- Drop table
DROP TABLE IF EXISTS leaks;