-- Drop the Triger first
DROP TRIGGER IF EXISTS update_events_updated_at ON events;

-- Drop the index
DROP INDEX IF EXISTS idx_events_tenant_id;
DROP INDEX IF EXISTS idx_events_provider_id;
DROP INDEX IF EXISTS idx_events_event_type;
DROP INDEX IF EXISTS idx_events_event_id;
DROP INDEX IF EXISTS idx_events_event_status;
DROP INDEX IF EXISTS idx_events_created_at;
DROP INDEX IF EXISTS idx_events_updated_at;

-- Drop the table
DROP TABLE IF EXISTS events;