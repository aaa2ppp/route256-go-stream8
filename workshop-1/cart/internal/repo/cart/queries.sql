-- name: List :many
SELECT * FROM cart WHERE user_id = $1;

-- name: Add :exec
INSERT INTO cart (user_id, sku, count) VALUES ($1, $2, $3)
ON CONFLICT (user_id, sku) DO UPDATE SET count = $3;

-- name: Delete :exec
DELETE FROM cart WHERE user_id = $1 AND sku = $2; 

-- name: Clear :exec
DELETE FROM cart WHERE user_id = $1;
