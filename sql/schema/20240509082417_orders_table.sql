-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
  id SERIAL NOT NULL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  total DECIMAL(10, 2) NOT NULL,
  status VARCHAR(20) CHECK (status IN ('pending', 'completed', 'cancelled')) NOT NULL DEFAULT 'pending',
  address TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders
-- +goose StatementEnd
