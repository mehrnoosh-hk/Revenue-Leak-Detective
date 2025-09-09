-- Create leak_type ENUM
CREATE TYPE leak_type_enum AS ENUM (
    'failed_payments',
    'unbilled_usage',
    'quiet_churn',
    'coupon_discount_misuse',
    'trial_forever',
    'other'
);

-- Create Leaks table
CREATE TABLE leaks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    leak_type leak_type_enum NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),  -- Money with 2 decimal places, must be positive
    confidence INTEGER NOT NULL CHECK (confidence >= 0 AND confidence <= 100),  -- 0-100 percentage
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- Create indexes for efficient lookups
CREATE INDEX idx_leaks_tenant_id ON leaks(tenant_id);
CREATE INDEX idx_leaks_customer_id ON leaks(customer_id);
CREATE INDEX idx_leaks_type ON leaks(leak_type);
CREATE INDEX idx_leaks_created_at ON leaks(created_at);

-- Create updated_at trigger for leaks
CREATE TRIGGER update_leaks_updated_at BEFORE UPDATE ON leaks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();