-- +goose Up
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    email TEXT NOT NULL UNIQUE,
    name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS property (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id),
    name TEXT NOT NULL,
    rent_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_property_user_id ON property (user_id);
-- note: indexes in your FKs makes JOINs faster, always index

CREATE TABLE IF NOT EXISTS tenant (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lease (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    property_id UUID NOT NULL UNIQUE,
    CONSTRAINT fk_property FOREIGN KEY (property_id) REFERENCES property (id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL,
    CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id),
    expected_rent_day SMALLINT NOT NULL DEFAULT 1,
    start_date TIMESTAMP NOT NULL DEFAULT now(),
    expiry_date TIMESTAMP,
    is_month_advance BOOLEAN NOT NULL DEFAULT FALSE,
    deposit_amount NUMERIC(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_lease_tenant_id ON lease (tenant_id);

CREATE INDEX idx_lease_property_id ON lease (property_id);

CREATE TABLE IF NOT EXISTS trade (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    lease_id UUID NOT NULL,
    CONSTRAINT fk_lease FOREIGN KEY (lease_id) REFERENCES lease (id) ON DELETE CASCADE,
    type SMALLINT NOT NULL, -- 0 = 'rent', 1 = 'deposit'
    paid_amount NUMERIC(10, 2) NOT NULL,
    start_date TIMESTAMP, -- [null] [rent] for type='rent', otherwise NULL
    end_date TIMESTAMP, -- [rent]
    note TEXT, --[null] [deposit] can be used for type='deposit', reason for reducing deposit
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- index to get newest payments from a specific tenant
CREATE INDEX idx_trade_lease_id_date_newest ON trade (lease_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS tenant CASCADE;

DROP TABLE IF EXISTS "user" CASCADE;

DROP TABLE IF EXISTS property CASCADE;

DROP TABLE IF EXISTS lease CASCADE;

DROP TABLE IF EXISTS trade CASCADE;