-- name: GetBySKU :one
SELECT * FROM stock
WHERE sku = $1;

-- name: Reserve :batchone
UPDATE stock SET reserved = $2
WHERE sku = $1 AND count-reserved >= $2
RETURNING sku;

-- name: ReserveRemove :batchone
UPDATE stock SET (count, reserved) = (count-$2, reserved-$2)
WHERE sku = $1 AND count >= $2 AND reserved >= $2
RETURNING sku;

-- name: ReserveCancel :batchone
UPDATE stock SET reserved = reserved-$2
WHERE sku = $1 AND reserved >= $2
RETURNING sku;
