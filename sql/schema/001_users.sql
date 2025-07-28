-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE users (
    id text NOT NULL PRIMARY KEY,
    name text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    is_admin integer DEFAULT 0,
    UNIQUE(name)
);

INSERT INTO users (
    id, name, password_hash, created_at, updated_at, is_admin
)
VALUES (
    '4d2ab25a-c902-4015-8ead-09f0e844d42e', 'Admin', '$2a$10$2WP4ssk27deQ7hGKWhwYl.DUlN740Gc0jDQwUT7eHIR8qXUcKnCw2', datetime('now'), datetime('now'), 1
);

-- +goose Down
DROP TABLE users;
