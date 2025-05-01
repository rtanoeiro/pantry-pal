-- name: AddItem :one
INSERT INTO pantry (
    id, user_id, item_name, quantity, added_at, expiry_at
) VALUES (
    ?, ?, ?, ?, strftime('%Y-%m-%d','now'), ?
)
RETURNING *;


-- UpdateItemQuantity :one
-- What'll see in the UI is a list of items, so we can probably use ID
UPDATE pantry
SET
    quantity = ?
WHERE id = ?

RETURNING *;

-- RemoveItem :one
DELETE FROM pantry
WHERE id = ?
RETURNING *;

-- GetAllItemsByName :many
-- Remember to lower the input from the UI
SELECT id, user_id, name, quantity, added_at, expiry_at
FROM pantry
WHERE user_id = ?
    AND lower(name) LIKE '%' || ? || '%'
ORDER BY added_at DESC;

-- GetOneItemByName :one
-- Remember to lower the input from the UI
SELECT id, user_id, name, quantity, added_at, expiry_at
FROM pantry
WHERE user_id = ?
    AND lower(name) = ?
