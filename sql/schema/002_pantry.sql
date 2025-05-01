-- +goose Up
CREATE TABLE pantry (
    id text NOT NULL PRIMARY KEY,
    user_id text NOT NULL,
    item_name text NOT NULL,
    quantity integer NOT NULL,
    added_at text NOT NULL,
    expiry_at text NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE pantry;