-- name: AddItemShopping :one
INSERT INTO cart_items (item_id, cart_id, item_name, quantity, added_at)
VALUES (?, ?, ?, ?, ?)

RETURNING *;

-- name: RemoveItemShopping :exec
DELETE FROM cart_items
WHERE item_id = ?;

-- name: GetAllShopping :many
SELECT 
    cart_items.cart_id, 
    cart_items.item_name, 
    cart_items.quantity, 
    cart_items.added_at 
FROM cart_items
INNER JOIN shopping_cart
ON shopping_cart.id = cart_items.cart_id
WHERE shopping_cart.user_id = ?;