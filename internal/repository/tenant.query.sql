-- name: CreateTenant :one
INSERT INTO
    tenant (
        id,
        lease_id,
        email,
        name,
        phone_number
    )
VALUES (
        gen_random_uuid (),
        $1,
        $2,
        $3,
        $4
    )
RETURNING
    *;

-- name: GetTenantById :one
SELECT
    id,
    lease_id,
    name,
    email,
    phone_number,
    created_at,
    updated_at
FROM tenant t
WHERE
    t.id = $1
LIMIT 1;

-- name: DeleteTenant :exec
DELETE FROM tenant WHERE id = $1;