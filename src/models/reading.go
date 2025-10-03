package models

import "time"

// Reading represents a flattened meter reading matching p1.meter_readings
type Reading struct {
	ID                    int64     `db:"id"`
	UniqueID              string    `db:"unique_id"`
	CreatedAt             time.Time `db:"created_at"`
	WifiSSID              string    `db:"wifi_ssid"`
	WifiStrength          int       `db:"wifi_strength"`
	SmrVersion            int       `db:"smr_version"`
	MeterModel            string    `db:"meter_model"`
	ActiveTariff          int       `db:"active_tariff"`
	TotalPowerImportKwh   float64   `db:"total_power_import_kwh"`
	TotalPowerImportT1Kwh float64   `db:"total_power_import_t1_kwh"`
	TotalPowerImportT2Kwh float64   `db:"total_power_import_t2_kwh"`
	TotalPowerExportKwh   float64   `db:"total_power_export_kwh"`
	TotalPowerExportT1Kwh float64   `db:"total_power_export_t1_kwh"`
	TotalPowerExportT2Kwh float64   `db:"total_power_export_t2_kwh"`
	ActivePowerW          float64   `db:"active_power_w"`
	ActivePowerL1W        float64   `db:"active_power_l1_w"`
	ActivePowerL2W        float64   `db:"active_power_l2_w"`
	ActivePowerL3W        float64   `db:"active_power_l3_w"`
	ActiveVoltageL1V      float64   `db:"active_voltage_l1_v"`
	ActiveVoltageL2V      float64   `db:"active_voltage_l2_v"`
	ActiveVoltageL3V      float64   `db:"active_voltage_l3_v"`
	ActiveCurrentA        float64   `db:"active_current_a"`
	ActiveCurrentL1A      float64   `db:"active_current_l1_a"`
	ActiveCurrentL2A      float64   `db:"active_current_l2_a"`
	ActiveCurrentL3A      float64   `db:"active_current_l3_a"`
	VoltageSagL1Count     int       `db:"voltage_sag_l1_count"`
	VoltageSagL2Count     int       `db:"voltage_sag_l2_count"`
	VoltageSagL3Count     int       `db:"voltage_sag_l3_count"`
	VoltageSwellL1Count   int       `db:"voltage_swell_l1_count"`
	VoltageSwellL2Count   int       `db:"voltage_swell_l2_count"`
	VoltageSwellL3Count   int       `db:"voltage_swell_l3_count"`
	AnyPowerFailCount     int       `db:"any_power_fail_count"`
	LongPowerFailCount    int       `db:"long_power_fail_count"`
	TotalGasM3            float64   `db:"total_gas_m3"`
	GasTimestamp          int64     `db:"gas_timestamp"`
	GasUniqueID           string    `db:"gas_unique_id"`
}

// ExternalReading represents entries in p1.external_readings
type ExternalReading struct {
	ID                   int64     `db:"id"`
	MeterReadingUniqueID string    `db:"meter_reading_unique_id"`
	CreatedAt            time.Time `db:"created_at"`
	UniqueID             string    `db:"unique_id"`
	Type                 string    `db:"type"`
	Timestamp            int64     `db:"timestamp"`
	Value                float64   `db:"value"`
	Unit                 string    `db:"unit"`
}
