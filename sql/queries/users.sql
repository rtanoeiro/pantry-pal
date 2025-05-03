-- name: CreateUser :one
INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)

RETURNING *;

-- name: UpdateUserEmail :exec
UPDATE users
SET 
    email = ?
WHERE id = ?
RETURNING *;

-- name: UpdateUserName :exec
UPDATE users
SET 
    name = ?
WHERE id = ?
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET 
    password_hash = ?
WHERE id = ?

RETURNING *;

-- name: UpdateUserAdmin :exec
UPDATE users
SET 
    is_admin = 1
WHERE id = ?

RETURNING *;

-- name: GetUserByIdOrEmail :one
SELECT id, name, email, password_hash, created_at, updated_at, is_admin
FROM users
WHERE id = ?1
   OR email = ?1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at, updated_at, is_admin
FROM users
WHERE email = ?;

-- name: ResetTable :exec
DELETE FROM users;
