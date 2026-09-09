-- name: CreateTenant :one
INSERT INTO
    tenant (id, email, name, phone_number)
VALUES (
        gen_random_uuid (),
        $1,
        $2,
        $3
    )
RETURNING
    *;