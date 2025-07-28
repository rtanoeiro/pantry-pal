-- +goose Up
CREATE TABLE cart_items (
    cart_id text NOT NULL PRIMARY KEY,
    item_name text NOT NULL,
    quantity INTEGER not null,
    added_at text NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES shopping_cart(id) on DELETE CASCADE
);