-- Create status enum
CREATE TYPE payment_status_enum AS ENUM (
    'pending',
    'succeeded',
    'failed',
    'other'
);

-- Create payment_type enum
CREATE TYPE payment_type_enum AS ENUM (
    'webhook',
    'history' -- When a tenant is onboarded, we need to fetch the history of the payments
);

-- Create payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    external_id VARCHAR(255) NOT NULL, -- Stripe payment intent ID
    amount DECIMAL(15,2) NOT NULL, -- The amount of the possible leak
    currency VARCHAR(3) NOT NULL,
    status payment_status_enum NOT NULL,
    payment_type payment_type_enum NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- Create indexes for efficient lookups
CREATE INDEX idx_payments_tenant_id ON payments(tenant_id);
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_external_id ON payments(external_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_payment_type ON payments(payment_type);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- Create updated_at trigger for payments
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();