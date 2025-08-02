-- name: AddItemShopping :exec
INSERT INTO cart_items (user_id, item_name, quantity)
VALUES (?, ?, ?);

-- name: UpdateItemShopping :exec
UPDATE cart_items
SET
    quantity = ?
WHERE item_name = ?
AND user_id = ?;

-- name: RemoveItemShopping :exec
DELETE FROM cart_items
WHERE item_name = ?
AND user_id = ?;

-- name: GetAllShopping :many
SELECT
    user_id,
    item_name, 
    quantity
FROM cart_items
WHERE user_id = ?;

-- name: FindItemShopping :one
SELECT
    item_name,
    quantity
FROM cart_items
WHERE item_name = ?
AND user_id = ?;