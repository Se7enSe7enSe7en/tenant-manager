-- +goose Up
CREATE TABLE identity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    provider TEXT NOT NULL, -- 'local' | 'google'
    provider_user_id TEXT NOT NULL, -- email for local, Google's `sub` for google
    password_hash TEXT, -- NULL unless provider='local'
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (user_id, provider_user_id)
);

CREATE INDEX idx_identity_user_id ON identity (user_id);

CREATE TABLE session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES "user" (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_session_user_id ON session (user_id);

-- +goose Down
DROP TABLE identity;

DROP TABLE session CASCADE;