CREATE TABLE orders
(
    id BIGSERIAL,
    number BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR NOT NULL,
    accrual FLOAT NOT NULL DEFAULT 0.0,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS orders_user_idx ON orders (user_id);
CREATE INDEX IF NOT EXISTS orders_number_idx ON orders (number);
CREATE INDEX IF NOT EXISTS orders_status_idx ON orders (status);