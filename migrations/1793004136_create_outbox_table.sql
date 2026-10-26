-- +goose Up
CREATE TABLE ticketing.outbox(
    id            UUID PRIMARY KEY DEFAULT uuidv7(),
    aggregate_id  UUID NOT NULL,             -- ticket id
    event_type    TEXT NOT NULL,             -- 'ticket.created'
    routing_key   TEXT NOT NULL,             -- service_code
    payload       JSONB NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at  TIMESTAMPTZ,
    attempts      INT NOT NULL DEFAULT 0,
    last_error    TEXT
);
CREATE INDEX outbox_unpublished_idx ON ticketing.outbox (created_at) WHERE published_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS ticketing.outbox;
