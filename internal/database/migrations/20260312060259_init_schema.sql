-- +goose Up
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS property (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id),
    name TEXT NOT NULL,
    rent_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tenant (
    id UUID PRIMARY KEY,
    property_id UUID NOT NULL,
    CONSTRAINT fk_property FOREIGN KEY (property_id) REFERENCES property (id),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    expected_rent_day SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS trade (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenant (id),
    property_id UUID NOT NULL,
    CONSTRAINT fk_property FOREIGN KEY (property_id) REFERENCES property (id),
    user_id UUID NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "user" (id),
    paid_amount NUMERIC(10, 2) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS tenant;

DROP TABLE IF EXISTS "user";

DROP TABLE IF EXISTS property;

DROP TABLE IF EXISTS trade;