// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: order_items.sql

package db

import (
	"context"
)

const createOrderItems = `-- name: CreateOrderItems :one
INSERT INTO order_items (order_id,product_id,quantity,price)
VALUES ($1,$2,$3,$4)
RETURNING id, order_id, product_id, quantity, price
`

type CreateOrderItemsParams struct {
	OrderID   int32
	ProductID int32
	Quantity  int32
	Price     string
}

func (q *Queries) CreateOrderItems(ctx context.Context, arg CreateOrderItemsParams) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, createOrderItems,
		arg.OrderID,
		arg.ProductID,
		arg.Quantity,
		arg.Price,
	)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.Price,
	)
	return i, err
}
