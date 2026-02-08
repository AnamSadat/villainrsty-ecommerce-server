-- name: GetUserByEmail :one
SELECT id, email, password, name, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (id, email, password, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUser :exec
UPDATE users
SET email = $2, password = $3, name = $4, updated_at = $5
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);