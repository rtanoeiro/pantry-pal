-- +goose Up
CREATE TABLE cart_items (
    user_id text NOT NULL,
    item_name text NOT NULL,
    quantity INTEGER not null,
    FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE
);

INSERT INTO cart_items VALUES ('4d2ab25a-c902-4015-8ead-09f0e844d42e', 'Rice', 2);

-- +goose Down
DROP TABLE cart_items;