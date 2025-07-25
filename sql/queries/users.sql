-- name: CreateUser :one
INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)

RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

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

-- name: MakeUserAdmin :exec
UPDATE users
SET 
    is_admin = 1
WHERE id = ?

RETURNING *;

-- name: RemoveUserAdmin :exec
UPDATE users
SET 
    is_admin = 0
WHERE id = ?
RETURNING *;

-- name: GetUserById :one
SELECT id, name, email, password_hash, created_at, updated_at, is_admin
FROM users
WHERE id = ?;

-- name: GetUserByName :one
SELECT id, name, email, password_hash, created_at, updated_at, is_admin
FROM users
WHERE name = ?;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at, updated_at, is_admin
FROM users
WHERE email = ?;

-- name: GetAllUsers :many
SELECT id, name, email, password_hash, created_at, updated_at, is_admin
FROM users
WHERE id != ?
ORDER BY created_at DESC;
