-- +goose Up
CREATE TABLE shopping_cart (
    id text NOT NULL PRIMARY KEY,
    user_id text NOT NULL,
    created_at timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE
);

INSERT INTO shopping_cart VALUES ('1', '4d2ab25a-c902-4015-8ead-09f0e844d42e', datetime('now'));

-- +goose Down
DROP TABLE shopping_cart;
