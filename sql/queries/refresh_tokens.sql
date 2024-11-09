-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens where token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET updated_at = $1, revoked_at = $1 where token = $2;

-- name: GetUserByRefreshToken :one
SELECT refresh_tokens.user_id from refresh_tokens INNER JOIN users on users.id = refresh_tokens.user_id where token = $1 AND revoked_at IS NULL AND expires_at > NOW();