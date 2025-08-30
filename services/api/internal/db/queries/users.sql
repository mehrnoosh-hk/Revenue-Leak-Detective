-- users table queries

-- name: GetUserByEmail :one
SELECT id, email, name, created_at, updated_at FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id, email, name, created_at, updated_at;

-- name: GetAllUsers :many
SELECT id, email, name, created_at, updated_at FROM users;

-- name: GetUserById :one
SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET email = $2, name = $3 WHERE id = $1 RETURNING id, email, name, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;