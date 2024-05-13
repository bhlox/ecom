// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: orders.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (user_id,total,address)
VALUES ($1,$2,$3)
RETURNING id,total
`

type CreateOrderParams struct {
	UserID  uuid.UUID
	Total   string
	Address string
}

type CreateOrderRow struct {
	ID    int32
	Total string
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (CreateOrderRow, error) {
	row := q.db.QueryRowContext(ctx, createOrder, arg.UserID, arg.Total, arg.Address)
	var i CreateOrderRow
	err := row.Scan(&i.ID, &i.Total)
	return i, err
}

const getOrder = `-- name: GetOrder :many
SELECT 
    orders.id, orders.user_id, orders.total, orders.status, orders.address, orders.created_at,order_items.id, order_items.order_id, order_items.product_id, order_items.quantity, order_items.price
FROM 
    orders
JOIN 
    order_items ON orders.id = order_items.order_id
WHERE 
    orders.id = $1
`

type GetOrderRow struct {
	Order     Order
	OrderItem OrderItem
}

func (q *Queries) GetOrder(ctx context.Context, id int32) ([]GetOrderRow, error) {
	rows, err := q.db.QueryContext(ctx, getOrder, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderRow
	for rows.Next() {
		var i GetOrderRow
		if err := rows.Scan(
			&i.Order.ID,
			&i.Order.UserID,
			&i.Order.Total,
			&i.Order.Status,
			&i.Order.Address,
			&i.Order.CreatedAt,
			&i.OrderItem.ID,
			&i.OrderItem.OrderID,
			&i.OrderItem.ProductID,
			&i.OrderItem.Quantity,
			&i.OrderItem.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
