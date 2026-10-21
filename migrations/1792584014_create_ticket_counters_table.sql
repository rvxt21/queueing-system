-- +goose Up
-- One row per service per day. Incremented atomically to hand out ticket numbers,
-- so two customers at the same moment never get the same number.
CREATE TABLE IF NOT EXISTS ticketing.ticket_counters (
    service_code TEXT NOT NULL REFERENCES ticketing.services (service_code) ON UPDATE CASCADE,
    day          DATE NOT NULL,
    last_number  INT  NOT NULL,
    PRIMARY KEY (service_code, day)
);

-- +goose Down
DROP TABLE IF EXISTS ticketing.ticket_counters;
