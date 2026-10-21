-- +goose Up
ALTER TABLE ticketing.tickets RENAME COLUMN service TO service_code;

ALTER TABLE ticketing.tickets
    ADD CONSTRAINT tickets_service_code_fkey
    FOREIGN KEY (service_code) REFERENCES ticketing.services (service_code)
    ON UPDATE CASCADE;

ALTER TABLE ticketing.tickets
    ADD CONSTRAINT tickets_status_check
    CHECK (status IN ('new', 'sent', 'delivered', 'resolved'));

-- +goose Down
ALTER TABLE ticketing.tickets DROP CONSTRAINT IF EXISTS tickets_status_check;
ALTER TABLE ticketing.tickets DROP CONSTRAINT IF EXISTS tickets_service_code_fkey;
ALTER TABLE ticketing.tickets RENAME COLUMN service_code TO service;
