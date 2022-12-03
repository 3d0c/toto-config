-- +goose Up
-- Please note, that regarding SQLite Type Affinity rules, any string datatype will
-- be converted into TEXT. Also, there is no requrements for field limitation provided.
-- This is why all string fields defined as type TEXT
CREATE TABLE toto_configuration (package TEXT NOT NULL, country_code TEXT NOT NULL DEFAULT 'ZZ', percentile_min INT DEFAULT 0, percentile_max INT DEFAULT 0, main_sku TEXT NOT NULL);

-- @TODO. Check whether is it possible to add integer rangine to the index
CREATE INDEX toto_cfg_name_cc_idx ON toto_configuration (package, country_code);

-- +goose Down
DROP TABLE toto_configuration;
