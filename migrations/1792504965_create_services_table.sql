-- +goose Up
CREATE TABLE IF NOT EXISTS ticketing.services (
    id            UUID PRIMARY KEY DEFAULT uuidv7(),
    service_code  TEXT NOT NULL UNIQUE,   -- 'passport', 'tax_code' (the contract with consumers)
    name          TEXT NOT NULL,          -- 'New passport' (shown in the frontend)
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS ticketing.services;
