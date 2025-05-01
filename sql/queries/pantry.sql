-- name: AddItem :one
INSERT INTO pantry (
    id, user_id, item_name, quantity, added_at, expiry_at
) VALUES (
    ?, ?, ?, ?, strftime('%Y-%m-%d','now'), ?
)
RETURNING *;


-- name: UpdateItemQuantity :one
-- What'll see in the UI is a list of items, so we can probably use ID
UPDATE pantry
SET
    quantity = ?
WHERE id = ?
    AND user_id = ?

RETURNING *;

-- name: RemoveItem :one
DELETE FROM pantry
WHERE id = ?
RETURNING *;

-- name: FindAllItemsByName :many
-- Remember to lower the input from the UI
SELECT id, user_id, item_name, quantity, added_at, expiry_at
FROM pantry
WHERE user_id = ?
    AND lower(item_name) LIKE '%' || ? || '%'
ORDER BY added_at DESC;

-- name: FindItemByName :many
-- Remember to lower the input from the UI
SELECT id, user_id, item_name, quantity, added_at, expiry_at
FROM pantry
WHERE user_id = ?
    AND lower(item_name) = ?
