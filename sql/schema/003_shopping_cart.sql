-- +goose Up
CREATE TABLE shopping_cart (
    id text NOT NULL PRIMARY KEY,
    user_id text NOT NULL,
    created_at text NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE
);

-- +goose Down
DROP TABLE shopping_cart;
