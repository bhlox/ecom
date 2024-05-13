-- name: CreateOrderItems :one
INSERT INTO order_items (order_id,product_id,quantity,price)
VALUES ($1,$2,$3,$4)
RETURNING *;