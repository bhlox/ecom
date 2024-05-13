-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS products_id_seq;

CREATE OR REPLACE FUNCTION generate_product_id() RETURNS TRIGGER AS $$
BEGIN
  NEW.id := nextval('products_id_seq');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS products (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  image VARCHAR(255) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  quantity INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_product_id
BEFORE INSERT ON products
FOR EACH ROW EXECUTE PROCEDURE generate_product_id();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products
-- +goose StatementEnd
