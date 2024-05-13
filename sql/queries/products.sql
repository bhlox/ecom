-- name: CreateProduct :one
INSERT INTO products (name,description,image,price,quantity)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetAllProducts :many
SELECT * FROM products;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1;

-- name: GetProductPriceAndStock :one
SELECT price,quantity FROM products
WHERE id = $1;

-- name: UpdateProductStock :exec
UPDATE products 
SET quantity = $1
WHERE id = $2;