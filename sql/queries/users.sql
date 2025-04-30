-- name: createUser :one
INSERT INTO users (id, name, email, password_hash) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: getUserById :one
SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id = $1

RETURNING id, name, email, password_hash, created_at, updated_at;

-- name: updateUserEmail :exect

UPDATE users
SET 
    EMAIL = $2
WHERE id = $1

RETURNING id, name, email, password_hash, created_at, updated_at;
