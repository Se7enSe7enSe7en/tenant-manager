-- name: CreateIdentity :one
INSERT INTO
    identity (
        id,
        user_id,
        provider,
        provider_user_id,
        password_hash
    )
VALUES (
        gen_random_uuid (),
        $1,
        $2,
        $3,
        $4
    ) RETURNING *;

-- name: GetIdentityByProvider :one
-- params: provider, provider_user_id
SELECT *
FROM identity
WHERE
    provider = $1
    AND provider_user_id = $2;

-- name: GetLocalIdentityByUserID :one
-- params: user_id
SELECT * FROM identity WHERE user_id = $1 AND provider = 'local';