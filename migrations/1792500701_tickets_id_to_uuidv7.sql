-- +goose Up
ALTER TABLE ticketing.tickets ALTER COLUMN id DROP DEFAULT;
ALTER TABLE ticketing.tickets ALTER COLUMN id TYPE UUID USING uuidv7();
ALTER TABLE ticketing.tickets ALTER COLUMN id SET DEFAULT uuidv7();
DROP SEQUENCE IF EXISTS ticketing.tickets_id_seq;

-- +goose Down
ALTER TABLE ticketing.tickets ALTER COLUMN id DROP DEFAULT;
CREATE SEQUENCE ticketing.tickets_id_seq;
ALTER TABLE ticketing.tickets ALTER COLUMN id TYPE BIGINT USING nextval('ticketing.tickets_id_seq');
ALTER TABLE ticketing.tickets ALTER COLUMN id SET DEFAULT nextval('ticketing.tickets_id_seq');
ALTER SEQUENCE ticketing.tickets_id_seq OWNED BY ticketing.tickets.id;
