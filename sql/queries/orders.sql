-- name: CreateOrder :one
INSERT INTO orders (user_id,total,address)
VALUES ($1,$2,$3)
RETURNING id,total;


-- name: GetOrder :many
SELECT 
    sqlc.embed(orders),sqlc.embed(order_items)
FROM 
    orders
JOIN 
    order_items ON orders.id = order_items.order_id
WHERE 
    orders.id = $1;
