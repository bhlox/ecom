-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items (
  id SERIAL PRIMARY KEY NOT NULL,
  order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  quantity INTEGER NOT NULL,
  price DECIMAL(10, 2) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items
-- +goose StatementEnd
