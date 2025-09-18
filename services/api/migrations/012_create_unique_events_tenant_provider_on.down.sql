-- Drop the unique composite index for (tenant_id, provider_id, event_id)
DROP INDEX IF EXISTS idx_events_tenant_provider_event_id;
