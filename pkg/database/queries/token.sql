
--UPDATE QUERIES
-- name: DiscardToken :exec
UPDATE tokens SET status='used', expires_at=NOW() WHERE uuid=$1;

-- name: BlacklistToken :exec
UPDATE tokens SET status='blacklisted', expires_at=NOW()  WHERE uuid=$1;

-- name: BlacklistTokenFamily :exec
UPDATE token_families SET status='blacklisted' WHERE uuid=$1;

---------------------------------------------
--INSERT QUERIES

-- name: CreateToken :one
INSERT INTO tokens(family, hash) VALUES($1, $2)
RETURNING *;

-- name: CreateTokenFamily :one
INSERT INTO token_families(user_uuid) VALUES($1)
RETURNING *;

---------------------------------------------
--GET QUERIES

-- name: GetTokenByHash :one
SELECT * FROM tokens WHERE hash = $1;

-- name: GetTokenFamilyByUUID :one
SELECT * FROM token_families WHERE uuid=$1;

---------------------------------------------
-- DELETE QUERIES

-- name: DeleteTokenByHash :exec
DELETE FROM tokens WHERE hash = $1;

-- name: DeleteTokensByFamily :exec
DELETE FROM tokens WHERE family = $1;


