-- name: CreateSession :one
INSERT INTO
    session (id, user_id, expires_at)
VALUES (gen_random_uuid (), $1, $2) RETURNING *;

-- name: GetSession :one
SELECT * FROM session WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM session WHERE id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM session WHERE now() >= expires_at;