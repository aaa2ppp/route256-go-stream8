-- name: Create :one
INSERT INTO "order" (user_id, status) VALUES ($1, 'new')
RETURNING order_id;

-- name: AddItems :exec
INSERT INTO order_items (order_id, sku, count)
SELECT UNNEST($1::bigint[]) AS order_id, UNNEST($2::bigint[]) AS sku, UNNEST($3::int[]) AS count;

-- name: GetByID :many
SELECT o.*, oi.sku, oi.count
FROM "order" AS o JOIN order_items AS oi USING(order_id)
WHERE o.order_id = $1;

-- name: SetStatus :exec
UPDATE "order" set status = $2
WHERE order_id = $1;
