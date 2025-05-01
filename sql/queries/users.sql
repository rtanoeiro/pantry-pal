-- name: createUser :one
INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)

RETURNING *;

-- name: updateUserEmail :exec

UPDATE users
SET 
    email = ?
WHERE id = ?

RETURNING *;

-- name: getUserById :one
SELECT id, name, email, password_hash, created_at, updated_at
FROM users
WHERE id = ?;