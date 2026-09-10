-- name: CreateLease :one
INSERT INTO
    lease (
        id,
        property_id,
        tenant_id,
        expected_rent_day,
        start_date,
        expiry_date,
        is_month_advance,
        deposit_amount,
        created_at,
        updated_at
    )
VALUES (
        gen_random_uuid (),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        now(),
        now()
    )
RETURNING
    *;

-- name: ListLease :many
SELECT
    -- lease
    l.id,
    l.expected_rent_day,
    l.start_date,
    l.expiry_date,
    l.is_month_advance,
    l.deposit_amount,
    l.created_at,
    -- property
    p.id AS property_id,
    p.name AS property_name,
    p.rent_amount AS property_rent_amount,
    -- tenant
    t.id AS tenant_id,
    t.email AS tenant_email,
    t.name AS tenant_name,
    t.phone_number AS tenant_phone_number,
    -- user
    u.id AS user_id,
    u.email AS user_email,
    u.name AS user_name
FROM
    lease l
    LEFT JOIN tenant t ON l.tenant_id = t.id
    LEFT JOIN property p ON l.property_id = p.id
    LEFT JOIN "user" u ON p.user_id = u.id
WHERE
    u.id = $1;

-- name: GetLeaseById :one
SELECT
    -- lease
    l.id,
    l.expected_rent_day,
    l.start_date,
    l.expiry_date,
    l.is_month_advance,
    l.deposit_amount,
    l.created_at,
    -- property
    p.id AS property_id,
    p.name AS property_name,
    p.rent_amount AS property_rent_amount,
    -- tenant
    t.id AS tenant_id,
    t.email AS tenant_email,
    t.name AS tenant_name,
    t.phone_number AS tenant_phone_number,
    -- user
    u.id AS user_id,
    u.email AS user_email,
    u.name AS user_name
FROM
    lease l
    LEFT JOIN tenant t ON l.tenant_id = t.id
    LEFT JOIN property p ON l.property_id = p.id
    LEFT JOIN "user" u ON p.user_id = u.id
WHERE
    l.id = $1
LIMIT 1;

-- name: UpdateLease :one
UPDATE lease
SET
    expected_rent_day = coalesce(
        sqlc.narg ('expected_rent_day'),
        expected_rent_day
    ),
    start_date = coalesce(
        sqlc.narg ('start_date'),
        start_date
    ),
    expiry_date = coalesce(
        sqlc.narg ('expiry_date'),
        expiry_date
    ),
    is_month_advance = coalesce(
        sqlc.narg ('is_month_advance'),
        is_month_advance
    ),
    deposit_amount = coalesce(
        sqlc.narg ('deposit_amount'),
        deposit_amount
    ),
    updated_at = now()
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteLease :exec
DELETE FROM lease WHERE id = $1;

-- name: CountLease :one
SELECT count(*)
FROM
    lease l
    LEFT JOIN property p ON property_id = p.id
    LEFT JOIN "user" u ON p.user_id = u.id
WHERE
    u.id = $1;

-- name: CountLeasePaid :one
SELECT count(*)
FROM
    lease l
    LEFT JOIN property p ON property_id = p.id
    LEFT JOIN "user" u ON p.user_id = u.id
WHERE
    u.id = $1
    AND now() < expiry_date;

-- name: CountLeaseUnpaid :one
SELECT count(*)
FROM
    lease l
    LEFT JOIN property p ON property_id = p.id
    LEFT JOIN "user" u ON p.user_id = u.id
WHERE
    u.id = $1
    AND expiry_date < now()
    AND now() < expiry_date + INTERVAL '1 month';

-- name: CountLeaseLate :one
SELECT count(*)
FROM
    lease l
    LEFT JOIN property p ON property_id = p.id
    LEFT JOIN "user" u ON p.user_id = u.id
WHERE
    u.id = $1
    AND expiry_date + INTERVAL '1 month' < now();