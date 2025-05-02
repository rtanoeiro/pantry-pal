-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE users (
    id text NOT NULL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP TRIGGER IF EXISTS remove_pantry_items;