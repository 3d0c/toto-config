-- +goose Up
-- Please note, that regarding SQLite Type Affinity rules, any string datatype will
-- be converted into TEXT. Also, there is no requrements for field limitation provided.
-- This is why all string fields defined as type TEXT
CREATE TABLE toto_configuration (package TEXT NOT NULL, country_code TEXT NOT NULL DEFAULT 'ZZ', percentile_min INT DEFAULT 0, percentile_max INT DEFAULT 0, main_SKU TEXT NOT NULL);

-- +goose Down
DROP TABLE toto_configuration;
