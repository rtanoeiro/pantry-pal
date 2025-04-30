-- +gooseUp
CREATE TABLE users (
    id uuid NOT NULL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

-- +gooseDown
DROP TABLE users;