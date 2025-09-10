-- Add the payment_id column to the leaks table
ALTER TABLE leaks ADD COLUMN payment_id uuid REFERENCES payments(id);

-- Add the foreign key constraint
ALTER TABLE leaks ADD CONSTRAINT fk_leaks_payment_id FOREIGN KEY (payment_id) REFERENCES payments(id);

-- Add the index
CREATE INDEX idx_leaks_payment_id ON leaks(payment_id);