-- name: CreateTrade :one
INSERT INTO
    trade (
        id,
        lease_id,
        type,
        paid_amount,
        start_date,
        end_date,
        note
    )
VALUES (
        gen_random_uuid (),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING
    *;