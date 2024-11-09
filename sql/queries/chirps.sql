-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps WHERE (user_id = $1 OR $1 IS NULL) ORDER BY created_at;

-- name: GetChirp :one
SELECT * FROM chirps where id = $1;

-- name: DeleteChirp :exec
DELETE from chirps where id = $1;