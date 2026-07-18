-- name: CreateUser :exec
INSERT INTO users (email, password_hash) VALUES($1, $2);

-- name: UpdateUserByUUID :exec
UPDATE users SET email=$2, password_hash=$3 WHERE uuid=$1;

-- name: GetUserByUUID :one
SELECT email from users WHERE uuid=$1;

-- name: DeleteUserByUUID :exec
DELETE FROM users WHERE uuid=$1;
