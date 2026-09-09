-- name: CreateProperty :one
INSERT INTO
    property (
        id,
        user_id,
        name,
        rent_amount
    )
VALUES (
        gen_random_uuid (),
        $1,
        $2,
        $3
    )
RETURNING
    *;

-- name: ListProperty :many
SELECT p.id, p.user_id, p.name, p.rent_amount, p.created_at, p.updated_at
FROM property p
WHERE
    p.user_id = $1;

-- name: ListUnoccupiedProperty :many
SELECT p.id, p.user_id, p.name, p.rent_amount, p.created_at, p.updated_at
FROM property p
    LEFT JOIN lease l ON p.id = l.property_id
WHERE
    l.property_id IS NULL
    AND p.user_id = $1;

-- name: GetPropertyById :one
SELECT p.id, p.user_id, p.rent_amount, p.created_at, p.updated_at
FROM property p
WHERE
    p.id = $1;