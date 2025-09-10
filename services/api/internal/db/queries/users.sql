-- users table queries

-- name: GetUserByEmail :one
SELECT id, email, name, external_id, created_at, updated_at FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, name, external_id)
VALUES (
    $1, 
    $2, 
    CASE WHEN sqlc.narg('external_id')::VARCHAR(255) IS NOT NULL THEN sqlc.narg('external_id')::VARCHAR(255) ELSE NULL END
)
RETURNING id, email, name, external_id, created_at, updated_at;

-- name: GetAllUsers :many
SELECT id, email, name, external_id, created_at, updated_at FROM users;

-- name: GetUserById :one
SELECT id, email, name, external_id, created_at, updated_at FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users 
SET 
    email = CASE WHEN sqlc.narg('email')::VARCHAR(255) IS NOT NULL THEN sqlc.narg('email')::VARCHAR(255) ELSE email END, 
    name = CASE WHEN sqlc.narg('name')::VARCHAR(255) IS NOT NULL THEN sqlc.narg('name')::VARCHAR(255) ELSE name END, 
    external_id = CASE WHEN sqlc.narg('external_id')::VARCHAR(255) IS NOT NULL THEN sqlc.narg('external_id')::VARCHAR(255) ELSE external_id END 
WHERE id = $1 
RETURNING id, email, name, external_id, created_at, updated_at;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = $1;