-- events table queries

-- name: GetEventByID :one
SELECT id, tenant_id, provider_id, event_type, event_id, status, data, created_at, updated_at FROM events WHERE id = $1;

-- name: CreateEvent :one
INSERT INTO events (tenant_id, provider_id, event_type, event_id, status, data) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, tenant_id, provider_id, event_type, event_id, status, data, created_at, updated_at;

-- name: GetAllEvents :many
SELECT id, tenant_id, provider_id, event_type, event_id, status, data, created_at, updated_at FROM events;

-- name: UpdateEvent :one
UPDATE events
SET
  tenant_id = CASE WHEN sqlc.narg('tenant_id')::uuid IS NOT NULL THEN sqlc.narg('tenant_id')::uuid ELSE tenant_id END,
  provider_id = CASE WHEN sqlc.narg('provider_id')::uuid IS NOT NULL THEN sqlc.narg('provider_id')::uuid ELSE provider_id END,
  event_type = CASE WHEN sqlc.narg('event_type')::event_type_enum IS NOT NULL THEN sqlc.narg('event_type')::event_type_enum ELSE event_type END,
  event_id = CASE WHEN sqlc.narg('event_id')::varchar IS NOT NULL THEN sqlc.narg('event_id')::varchar ELSE event_id END,
  status = CASE WHEN sqlc.narg('status')::event_status_enum IS NOT NULL THEN sqlc.narg('status')::event_status_enum ELSE status END,
  data = CASE WHEN sqlc.narg('data')::jsonb IS NOT NULL THEN sqlc.narg('data')::jsonb ELSE data END
WHERE id = sqlc.arg('id')
RETURNING id, tenant_id, provider_id, event_type, event_id, status, data, created_at, updated_at;

-- name: DeleteEvent :execrows
DELETE FROM events WHERE id = $1;