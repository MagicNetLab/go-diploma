CREATE TABLE withdraw
(
    id BIGSERIAL,
    order_num BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    sum FLOAT NOT NULL DEFAULT 0.0,
    processed_at timestamp DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS withdraw_user_idx ON withdraw (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS withdraw_order_idx ON withdraw (order_num);
