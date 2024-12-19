-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, passwd)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key=$1;

-- name: GetUserByName :one
SELECT * FROM users WHERE name=$1;
