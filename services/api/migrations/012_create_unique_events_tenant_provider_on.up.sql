-- Create a unique composite index to ensure no duplicate events per tenant and provider
-- Uniqueness: (tenant_id, provider_id, event_id)
CREATE UNIQUE INDEX idx_events_tenant_provider_event_id ON events(tenant_id, provider_id, event_id);
