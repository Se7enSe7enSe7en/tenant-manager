-- name: CreateTenant :one
INSERT INTO
    tenant (
        id,
        email,
        name,
        phone_number,
        expected_rent_day,
        property_id
    )
VALUES (
        gen_random_uuid (),
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING
    *;

-- name: ListTenantsWithProperty :many
SELECT
    t.id AS tenant_id,
    t.email AS tenant_email,
    t.name AS tenant_name,
    t.expected_rent_day AS tenant_expected_rent_day,
    t.phone_number AS tenant_phone_number,
    
    p.id AS property_id,
    p.name AS property_name,
    p.rent_amount AS property_rent_amount,
    p.updated_at AS property_updated_at
FROM tenant t
LEFT JOIN property p
    ON t.property_id = p.id
WHERE p.user_id = $1;

