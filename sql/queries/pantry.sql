-- addItem :one
INSERT INTO pantry (id, user_id, name, quantity, added_at, expiry_at)
VALUES ($1, $2, $3, $4, $5, $6);

RETURNING id, user_id, name, quantity, added_at, expiry_at;


-- updateItemQuantity :one
-- What'll see in the UI is a list of items, so we can probably use ID
UPDATE pantry
SET quantity = $1,
WHERE id = $2

RETURNING id, user_id, name, quantity, added_at, expiry_at;

-- removeItem :one
DELETE FROM pantry
WHERE id = $1
RETURNING id, user_id, name, added_at, expiry_at;

-- getAllItemsByName :many
-- Remember to lower the input from the UI
SELECT id, user_id, name, quantity, added_at, expiry_at
FROM pantry
WHERE user_id = $1
    AND lower(name) ILIKE '%' || $2 || '%'
ORDER BY added_at DESC;

-- getOneItemByName :one
-- Remember to lower the input from the UI
SELECT id, user_id, name, quantity, added_at, expiry_at
FROM pantry
WHERE user_id = $1
    AND lower(name) = $2

