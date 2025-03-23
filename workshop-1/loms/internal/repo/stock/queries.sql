-- name: GetBySKU :one
SELECT * FROM stock
WHERE sku = @sku;

-- name: Reserve :batchone
WITH rows AS (
    UPDATE stock SET available = available - @count, reserved = reserved + @count
    WHERE sku = @sku AND available >= @count
    RETURNING sku
)
SELECT COUNT(*);

-- name: ReserveCancel :batchone
WITH rows AS (
    UPDATE stock SET available = available + @count, reserved = reserved - @count 
    WHERE sku = @sku AND reserved >= @count
    RETURNING sku
)
SELECT COUNT(*);

-- name: ReserveRemove :batchone
WITH rows AS (
    UPDATE stock SET reserved = reserved - @count
    WHERE sku = @sku AND reserved >= @count
    RETURNING sku
)
SELECT COUNT(*);

