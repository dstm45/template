-- +goose Up
CREATE TYPE role AS ENUM('standard', 'admin');
CREATE TYPE token_status AS ENUM('active', 'used', 'blacklisted');
CREATE TYPE token_family_status AS ENUM('active', 'blacklisted');
CREATE TABLE IF NOT EXISTS users(
  uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  email TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  role role NOT NULL DEFAULT 'standard'::role
);


CREATE TABLE IF NOT EXISTS standards(
	user_uuid uuid NOT NULL UNIQUE REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS admins(
	user_uuid uuid NOT NULL UNIQUE REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_families(
  id BIGSERIAL PRIMARY KEY,
  uuid uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
  user_uuid uuid references users(uuid) NOT NULL,
  status token_family_status NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS tokens(
  id BIGSERIAL PRIMARY KEY,
  hash TEXT NOT NULL,
  uuid uuid NOT NULL DEFAULT gen_random_uuid(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMP WITH TIME ZONE,
  family uuid references token_families(uuid) NOT NULL,
  status token_status NOT NULL DEFAULT 'active'
);

CREATE VIEW user_public_data AS SELECT email, uuid, role FROM users;

-- +goose Down
DROP TABLE tokens;
DROP TABLE token_families;
DROP TABLE standards;
DROP TABLE admins;
DROP VIEW user_public_data;
DROP TABLE users;
DROP TYPE token_status;
DROP TYPE token_family_status;
DROP TYPE role;
