-- +goose Up
ALTER TABLE ticketing.services ADD COLUMN prefix TEXT;

UPDATE ticketing.services AS s SET prefix = v.prefix
FROM (VALUES
    ('passport_new',            'PN'),
    ('passport_collect',        'PC'),
    ('id_card_new',             'IN'),
    ('id_card_collect',         'IC'),
    ('tax_registration',        'TR'),
    ('tax_declaration',         'TD'),
    ('tax_certificate',         'TC'),
    ('address_registration',    'AR'),
    ('address_deregistration',  'AD'),
    ('vehicle_registration',    'VR'),
    ('driving_license_new',     'DN'),
    ('driving_license_replace', 'DR'),
    ('birth_certificate',       'BC'),
    ('marriage_certificate',    'MC'),
    ('marriage_registration',   'MR'),
    ('business_registration',   'BR'),
    ('criminal_record',         'CR'),
    ('parking_permit',          'PP'),
    ('building_permit',         'BP'),
    ('complaint',               'CO')
) AS v (service_code, prefix)
WHERE s.service_code = v.service_code;

ALTER TABLE ticketing.services
    ALTER COLUMN prefix SET NOT NULL,
    ADD CONSTRAINT services_prefix_key UNIQUE (prefix),
    ADD CONSTRAINT services_prefix_check CHECK (prefix ~ '^[A-Z]{1,3}$');

-- +goose Down
ALTER TABLE ticketing.services DROP COLUMN IF EXISTS prefix;
