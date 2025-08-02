-- +goose Up
CREATE TABLE cart_items (
    cart_id text NOT NULL PRIMARY KEY,
    item_id text NOT NULL,
    item_name text NOT NULL,
    quantity INTEGER not null,
    added_at timestamp NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES shopping_cart(id) on DELETE CASCADE
);

INSERT INTO cart_items VALUES ('1', '4d2ab25a-c902', 'Rice', 2, datetime('now'));

-- +goose Down
DROP TABLE cart_items;