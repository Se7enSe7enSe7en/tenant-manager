-- name: ListProperties :many
SELECT * FROM property;

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

-- name: ListUnoccupiedProperties :many
SELECT p.*
FROM property p
LEFT JOIN tenant t 
    ON p.id = t.property_id
WHERE t.property_id IS NULL
AND p.user_id = $1;