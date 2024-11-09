-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users where email = $1;

-- name: UpdateUser :one
UPDATE users 
SET 
    email = COALESCE($1, email),
    hashed_password = COALESCE($2, hashed_password),
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpdateUserRed :exec
UPDATE users 
SET 
    is_chirpy_red = true,
    updated_at = NOW()
WHERE id = $1;
