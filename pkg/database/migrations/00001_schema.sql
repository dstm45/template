-- +goose Up
CREATE TABLE IF NOT EXISTS users(
  uuid UUID NOT NULL DEFAULT gen_random_uuid(),
  email TEXT NOT NULL,
  password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
