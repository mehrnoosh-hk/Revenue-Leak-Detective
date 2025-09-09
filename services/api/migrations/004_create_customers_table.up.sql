-- Create customers table with tenant relationship and external_id
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    external_id VARCHAR(255) NOT NULL,  -- Customer ID within the payment integration
    email CITEXT NOT NULL, 
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE(tenant_id, external_id)  -- Prevents duplicate external_id per tenant
);

-- Create indexes for efficient lookups
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_tenant_id ON customers(tenant_id);
CREATE INDEX idx_customers_external_id ON customers(external_id);

-- Create updated_at trigger for customers
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();