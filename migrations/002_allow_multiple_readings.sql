-- Migration: allow multiple meter_readings per meter unique_id
-- 1) remove UNIQUE constraint on meter_readings.unique_id
-- 2) create new external_readings referencing meter_readings(id)
-- 3) migrate existing external_readings data and rename
BEGIN;

-- drop unique constraint if it exists (name typically meter_readings_unique_id_key)
ALTER TABLE
	IF EXISTS p1.meter_readings DROP CONSTRAINT IF EXISTS meter_readings_unique_id_key;

-- create new table for externals referencing the primary key id
CREATE TABLE IF NOT EXISTS p1.external_readings_new (
	id BIGSERIAL PRIMARY KEY,
	meter_reading_id BIGINT NOT NULL REFERENCES p1.meter_readings(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	unique_id TEXT,
	type TEXT,
	timestamp BIGINT,
	value NUMERIC(14, 3),
	unit TEXT
);

-- copy data from old externals if present (map by meter unique_id -> meter_readings.id)
INSERT INTO
	p1.external_readings_new (
		meter_reading_id,
		created_at,
		unique_id,
		type,
		timestamp,
		value,
		unit
	)
SELECT
	mr.id,
	er.created_at,
	er.unique_id,
	er.type,
	er.timestamp,
	er.value,
	er.unit
FROM
	p1.external_readings er
	JOIN p1.meter_readings mr ON mr.unique_id = er.meter_reading_unique_id;

-- drop old external table and rename the new one
DROP TABLE IF EXISTS p1.external_readings;

ALTER TABLE
	p1.external_readings_new RENAME TO external_readings;

COMMIT;