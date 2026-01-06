-- Drop columns from meter_readings table
ALTER TABLE p1.meter_readings DROP COLUMN IF EXISTS unique_id;
ALTER TABLE p1.meter_readings DROP COLUMN IF EXISTS wifi_ssid;
ALTER TABLE p1.meter_readings DROP COLUMN IF EXISTS wifi_strength;
ALTER TABLE p1.meter_readings DROP COLUMN IF EXISTS smr_version;
ALTER TABLE p1.meter_readings DROP COLUMN IF EXISTS meter_model;
ALTER TABLE p1.meter_readings DROP COLUMN IF EXISTS gas_unique_id;
