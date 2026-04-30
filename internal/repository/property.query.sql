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
    ) RETURNING *;