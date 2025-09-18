-- Drop the constraint
ALTER TABLE leaks DROP CONSTRAINT fk_leaks_payment_id;

-- Drop the column
ALTER TABLE leaks DROP COLUMN payment_id;

-- Drop the index
DROP INDEX IF EXISTS idx_leaks_payment_id;