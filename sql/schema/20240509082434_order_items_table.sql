-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS order_items_id_seq;

CREATE OR REPLACE FUNCTION generate_order_item_id() RETURNS TRIGGER AS $$
BEGIN
  NEW.id := nextval('order_items_id_seq');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS order_items (
  id INTEGER PRIMARY KEY NOT NULL,
  order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  quantity INTEGER NOT NULL,
  price DECIMAL(10, 2) NOT NULL
);

CREATE TRIGGER set_order_item_id
BEFORE INSERT ON order_items
FOR EACH ROW EXECUTE PROCEDURE generate_order_item_id();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items
-- +goose StatementEnd
