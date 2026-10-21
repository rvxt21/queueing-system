-- +goose Up
INSERT INTO ticketing.services (service_code, name) VALUES
    ('passport_new',            'Apply for a new passport'),
    ('passport_collect',        'Collect a passport'),
    ('id_card_new',             'Apply for a new ID card'),
    ('id_card_collect',         'Collect an ID card'),
    ('tax_registration',        'Tax registration or update'),
    ('tax_declaration',         'Submit a tax declaration'),
    ('tax_certificate',         'Certificate of no tax arrears'),
    ('address_registration',    'Register a place of residence'),
    ('address_deregistration',  'Deregister from a place of residence'),
    ('vehicle_registration',    'Register a vehicle'),
    ('driving_license_new',     'Apply for a driving license'),
    ('driving_license_replace', 'Replace a driving license'),
    ('birth_certificate',       'Request a birth certificate'),
    ('marriage_certificate',    'Request a marriage certificate'),
    ('marriage_registration',   'Marriage registration'),
    ('business_registration',   'Register a business'),
    ('criminal_record',         'Certificate of no criminal record'),
    ('parking_permit',          'Resident parking permit'),
    ('building_permit',         'Building permit application'),
    ('complaint',               'Submit a complaint or general inquiry')
ON CONFLICT (service_code) DO UPDATE
    SET name = EXCLUDED.name,
        updated_at = now();

-- +goose Down
-- Fails if tickets still reference these services (foreign key), which protects ticket data.
DELETE FROM ticketing.services WHERE service_code IN (
    'passport_new', 'passport_collect', 'id_card_new', 'id_card_collect',
    'tax_registration', 'tax_declaration', 'tax_certificate',
    'address_registration', 'address_deregistration', 'vehicle_registration',
    'driving_license_new', 'driving_license_replace',
    'birth_certificate', 'marriage_certificate', 'marriage_registration',
    'business_registration', 'criminal_record',
    'parking_permit', 'building_permit', 'complaint'
);
