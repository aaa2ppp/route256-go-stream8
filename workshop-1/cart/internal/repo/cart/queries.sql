-- name: List :many
SELECT user_id, sku, count FROM cart WHERE user_id = $1;

-- name: AddBatch :batchexec
INSERT INTO cart (user_id, sku, count) VALUES ($1, $2, $3)
ON CONFLICT (user_id, sku) DO UPDATE set count = count + $3;

-- name: AddArrays :exec
INSERT INTO cart (user_id, sku, count)
SELECT 
    UNNEST($1::bigint[]) AS user_id,
    UNNEST($2::int[]   ) AS sku,
    UNNEST($3::int[]   ) AS count    
ON CONFLICT (user_id, sku) 
DO UPDATE SET count = cart.count + EXCLUDED.count;

-- name: Delete :exec
DELETE FROM cart WHERE user_id = $1 AND sku = $2; 

-- name: Clear :exec
DELETE FROM cart WHERE user_id = $1;

