-- name: CreateUser :one
INSERT INTO users (
    id,
    lastname,
    firstname,
    email,
    phone,
    address,
    hashed_password
) VALUES (
    $1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListAllUsers :many
SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users SET hashed_password = $2, updated_at = $3 WHERE id = $1 RETURNING *;

-- name: UpdateUser :one
UPDATE users SET address = $4, phone = $3, email = $2, updated_at = $5 WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;