-- Create event_type ENUM
CREATE TYPE event_type_enum AS ENUM (
    'payment_failed',
    'payment_succeeded',
    'payment_refunded',
    'payment_updated'
);

-- Create event_status ENUM
CREATE TYPE event_status_enum AS ENUM (
    'pending',
    'processed',
    'failed'
);

-- Create events table, this table is used to store any webhook events of the system
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    provider_id UUID NOT NULL, -- Providers like Stripe, Chargebee, Recurly, etc.
    event_type event_type_enum NOT NULL, -- Event type like payment_failed, payment_succeeded, etc.
    event_id VARCHAR(255) NOT NULL, -- Id of the event from the provider
    status event_status_enum NOT NULL, -- Event status like based on it's status from the leak engine
    data JSONB NOT NULL, -- Event data like payment_id, amount, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (provider_id) REFERENCES providers(id) ON DELETE CASCADE
);

-- Create indexes for efficient lookups
CREATE INDEX idx_events_tenant_id ON events(tenant_id);
CREATE INDEX idx_events_provider_id ON events(provider_id);
CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_event_id ON events(event_id);
CREATE INDEX idx_events_event_status ON events(status);
CREATE INDEX idx_events_created_at ON events(created_at);
CREATE INDEX idx_events_updated_at ON events(updated_at);

-- Create updated_at trigger for events
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();