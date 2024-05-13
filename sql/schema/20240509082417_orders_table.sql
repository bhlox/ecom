-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS orders_id_seq;

CREATE OR REPLACE FUNCTION generate_order_id() RETURNS TRIGGER AS $$
BEGIN
  NEW.id := nextval('orders_id_seq');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS orders (
  id INTEGER NOT NULL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  total DECIMAL(10, 2) NOT NULL,
  status VARCHAR(20) CHECK (status IN ('pending', 'completed', 'cancelled')) NOT NULL DEFAULT 'pending',
  address TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_order_id
BEFORE INSERT ON orders
FOR EACH ROW EXECUTE PROCEDURE generate_order_id();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders
-- +goose StatementEnd
