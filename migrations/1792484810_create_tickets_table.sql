-- +goose Up
CREATE SCHEMA IF NOT EXISTS ticketing;
CREATE TABLE IF NOT EXISTS ticketing.tickets (
    id          BIGSERIAL PRIMARY KEY,
    number      TEXT NOT NULL,          -- e.g. "P-042"
    service     TEXT NOT NULL,          -- passport / tax_code / ...
    status      TEXT NOT NULL DEFAULT 'new',  -- new -> delivered -> resolved
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS ticketing.tickets;
DROP SCHEMA IF EXISTS ticketing;
