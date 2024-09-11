CREATE TABLE users
(
    id BIGSERIAL,
    login VARCHAR NOT NULL,
    password VARCHAR NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS login_uniq_idx ON users (login);