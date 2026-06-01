-- name: CreateUser :one
INSERT INTO
    "user" (id, email, name)
VALUES (gen_random_uuid (), $1, $2) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM "user" WHERE id = $1;