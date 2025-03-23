// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package queries

import (
	"context"
)

const addItems = `-- name: AddItems :exec
INSERT INTO order_items (order_id, sku, count)
SELECT UNNEST($1::bigint[]) AS order_id, UNNEST($2::bigint[]) AS sku, UNNEST($3::int[]) AS count
`

type AddItemsParams struct {
	Orders []int64 `json:"orders"`
	Skus   []int64 `json:"skus"`
	Counts []int32 `json:"counts"`
}

func (q *Queries) AddItems(ctx context.Context, arg AddItemsParams) error {
	_, err := q.db.Exec(ctx, addItems, arg.Orders, arg.Skus, arg.Counts)
	return err
}

const create = `-- name: Create :one
INSERT INTO "order" (user_id, status) VALUES ($1, 'new')
RETURNING order_id
`

func (q *Queries) Create(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRow(ctx, create, userID)
	var order_id int64
	err := row.Scan(&order_id)
	return order_id, err
}

const getByID = `-- name: GetByID :many
SELECT o.order_id, o.user_id, o.status, oi.sku, oi.count
FROM "order" AS o JOIN order_items AS oi USING(order_id)
WHERE o.order_id = $1
`

type GetByIDRow struct {
	OrderID int64       `json:"order_id"`
	UserID  int64       `json:"user_id"`
	Status  OrderStatus `json:"status"`
	Sku     int32       `json:"sku"`
	Count   int32       `json:"count"`
}

func (q *Queries) GetByID(ctx context.Context, orderID int64) ([]GetByIDRow, error) {
	rows, err := q.db.Query(ctx, getByID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetByIDRow
	for rows.Next() {
		var i GetByIDRow
		if err := rows.Scan(
			&i.OrderID,
			&i.UserID,
			&i.Status,
			&i.Sku,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setStatus = `-- name: SetStatus :exec
UPDATE "order" set status = $2
WHERE order_id = $1
`

type SetStatusParams struct {
	OrderID int64       `json:"order_id"`
	Status  OrderStatus `json:"status"`
}

func (q *Queries) SetStatus(ctx context.Context, arg SetStatusParams) error {
	_, err := q.db.Exec(ctx, setStatus, arg.OrderID, arg.Status)
	return err
}
