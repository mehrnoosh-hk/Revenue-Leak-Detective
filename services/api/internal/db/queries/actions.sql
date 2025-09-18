-- name: GetActionByID :one
SELECT id, leak_id, action_type, status, result, created_at, updated_at
FROM actions
WHERE id = $1;

-- name: CreateAction :one
INSERT INTO actions (leak_id, action_type, status, result) VALUES ($1, $2, $3, $4) RETURNING id, leak_id, action_type, status, result, created_at, updated_at;

-- name: GetAllActions :many
SELECT id, leak_id, action_type, status, result, created_at, updated_at
FROM actions;

-- name: GetAllActionsPaginated :many
SELECT id, leak_id, action_type, status, result, created_at, updated_at
FROM actions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllActions :one
SELECT COUNT(*) FROM actions;


-- name: UpdateAction :one
UPDATE actions 
SET 
    action_type = CASE WHEN sqlc.narg('action_type')::action_type_enum IS NOT NULL THEN sqlc.narg('action_type')::action_type_enum ELSE action_type END, 
    status = CASE WHEN sqlc.narg('status')::action_status_enum IS NOT NULL THEN sqlc.narg('status')::action_status_enum ELSE status END, 
    result = CASE WHEN sqlc.narg('result')::action_result_enum IS NOT NULL THEN sqlc.narg('result')::action_result_enum ELSE result END 
WHERE id = $1 
RETURNING id, leak_id, action_type, status, result, created_at, updated_at;

-- name: DeleteAction :execrows
DELETE FROM actions WHERE id = $1;