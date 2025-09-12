-- name: GetActionByIDForTenant :one
SELECT a.id, a.leak_id, a.action_type, a.status, a.result, a.created_at, a.updated_at
FROM actions a
JOIN leaks l ON l.id = a.leak_id
WHERE a.id = $1 AND l.tenant_id = $2;

-- name: CreateAction :one
INSERT INTO actions (leak_id, action_type, status, result) VALUES ($1, $2, $3, $4) RETURNING id, leak_id, action_type, status, result, created_at, updated_at;

-- name: GetAllActionsForTenant :many
SELECT a.id, a.leak_id, a.action_type, a.status, a.result, a.created_at, a.updated_at
FROM actions a
JOIN leaks l ON l.id = a.leak_id
WHERE l.tenant_id = $1;

-- name: UpdateAction :one
UPDATE actions 
SET 
    leak_id = CASE WHEN sqlc.narg('leak_id')::uuid IS NOT NULL THEN sqlc.narg('leak_id')::uuid ELSE leak_id END, 
    action_type = CASE WHEN sqlc.narg('action_type')::action_type_enum IS NOT NULL THEN sqlc.narg('action_type')::action_type_enum ELSE action_type END, 
    status = CASE WHEN sqlc.narg('status')::action_status_enum IS NOT NULL THEN sqlc.narg('status')::action_status_enum ELSE status END, 
    result = CASE WHEN sqlc.narg('result')::action_result_enum IS NOT NULL THEN sqlc.narg('result')::action_result_enum ELSE result END 
WHERE id = $1 
RETURNING id, leak_id, action_type, status, result, created_at, updated_at;

-- name: DeleteAction :execrows
DELETE FROM actions WHERE id = $1;